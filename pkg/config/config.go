package config

import (
	"context"
	"errors"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

type ConfigManager interface {
	GetData() interface{}                          // 获取配置数据
	GetErr() error                                 // 获取错误信息
	PutData()                                      // 本地配置文件不需要，空实现
	LoadData()                                     // 先加载再获取数据
	WatchData(ctx context.Context) <-chan struct{} //热更新
}

// yaml和json部分
type LocalConfig struct {
	Viper      *viper.Viper
	Lastupdate time.Time
	Data       interface{} // 传入的配置结构体指针
	Path       string      // 配置文件路径
	Err        error       // 统一在这里存放错误
}

// c需传入结构体指针，返回接口类型
func NewLocalConfig(c interface{}, path string) (ConfigManager, error) {
	if reflect.ValueOf(c).Kind() != reflect.Ptr {
		return nil, errors.New("config struct must be a pointer")
	}
	return &LocalConfig{
		Viper:      viper.New(),
		Data:       c,
		Lastupdate: time.Now(),
		Path:       path,
		Err:        nil,
	}, nil
}

func (l *LocalConfig) GetData() interface{} {
	return l.Data
}

func (l *LocalConfig) GetErr() error {
	return l.Err
}

func (l *LocalConfig) LoadData() {
	l.Viper.SetConfigFile(l.Path)

	// 检查扩展名(支持yaml, json)，如果不是则报错
	ext := filepath.Ext(l.Path)
	switch ext {
	case ".yaml", ".yml":
		l.Viper.SetConfigType("yaml")
	case ".json":
		l.Viper.SetConfigType("json")
	default:
		l.Err = errors.New("only .yaml, .yml, or .json are supported")
		return
	}

	// 读取配置文件
	err := loadData(l)
	l.Err = err

}

// put感觉本地用不到，所以空实现
func (l *LocalConfig) PutData() {}

// 传回的通道会在配置文件发生变化时发送信号
func (l *LocalConfig) WatchData(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{}, 10)

	go func() {
		defer close(ch)
		l.Viper.WatchConfig()
		l.Viper.OnConfigChange(func(e fsnotify.Event) {
			// 不知道为什么会连着调用两次，所以加个时间限制
			if time.Since(l.Lastupdate) < 1*time.Second {
				return
			}

			err := loadData(l)
			if err != nil {
				// 热更新时出错
				l.Err = err
				return
			}

			// 发送信号通知配置已更新
			select {
			case ch <- struct{}{}:
			case <-ctx.Done():
				// 如果接收到停止信号，则退出
				return
			}
		})

		<-ctx.Done()
	}()

	return ch
}

func loadData(l *LocalConfig) error {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	if err := l.Viper.ReadInConfig(); err != nil {
		return err
	}

	if err := l.Viper.Unmarshal(l.Data); err != nil {
		return err
	}

	l.Lastupdate = time.Now()
	return nil
}

// Nacos部分
// 定义查询结构体
type NacosConfig struct {
	DataId string                      // 配置key
	Group  string                      // 配置组
	Data   string                      // 配置value，获取值时也在这里获取
	Client config_client.IConfigClient // Nacos客户端
	Err    error                       // 统一在这里存放错误
}

func NewClientConfig(namespace, username, password string, time uint64) *constant.ClientConfig {
	return &constant.ClientConfig{
		NamespaceId: namespace,
		TimeoutMs:   time,
		Username:    username,
		Password:    password,
	}
}

func NewServerConfig(ip []string, port []uint64) []constant.ServerConfig {
	serverConfigs := make([]constant.ServerConfig, 0, len(ip))
	for i := 0; i < len(ip); i++ {
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr: ip[i],
			Port:   port[i],
		})
	}
	return serverConfigs
}

// 创建nacos客户端
func NewNacos(clientConfig *constant.ClientConfig, serverConfigs []constant.ServerConfig) (config_client.IConfigClient, error) {
	client, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

// 返回接口类型
func NewNacosConfig(dataId, group, content string, clientconfig *constant.ClientConfig, serverconfig []constant.ServerConfig) (ConfigManager, error) {
	client, err := NewNacos(clientconfig, serverconfig)
	if err != nil {
		return nil, err
	}
	return &NacosConfig{
		DataId: dataId,
		Group:  group,
		Data:   content,
		Client: client,
		Err:    nil,
	}, nil
}

func (p *NacosConfig) GetData() interface{} {
	return p.Data
}

func (p *NacosConfig) GetErr() error {
	return p.Err
}

func (p *NacosConfig) PutData() {
	_, err := p.Client.PublishConfig(vo.ConfigParam{
		DataId:  p.DataId,
		Group:   p.Group,
		Content: p.Data,
	})
	p.Err = err
}

func (p *NacosConfig) LoadData() {
	content, err := p.Client.GetConfig(vo.ConfigParam{
		DataId: p.DataId,
		Group:  p.Group,
	})
	if err != nil {
		p.Err = err
		return
	}
	p.Data = content
}

func (p *NacosConfig) WatchData(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{}, 10)
	go func() {
		defer close(ch)
		err := p.Client.ListenConfig(vo.ConfigParam{
			DataId: p.DataId,
			Group:  p.Group,
			OnChange: func(namespace, group, dataId, data string) {
				var mu sync.Mutex
				mu.Lock()
				defer mu.Unlock()
				p.Data = data
				ch <- struct{}{}
			},
		})
		if err != nil {
			var mu sync.Mutex
			mu.Lock()
			defer mu.Unlock()
			p.Err = err
			return
		}
		<-ctx.Done()
	}()
	return ch
}
