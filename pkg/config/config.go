package config

import (
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
	GetData() interface{}    // 获取配置数据
	LoadData() error         // 先加载再获取数据
	WatchData() <-chan error //热更新
	Close() error            //关闭
}

// yaml和json部分
type LocalConfig struct {
	Viper      *viper.Viper
	Lastupdate time.Time
	Data       interface{} // 传入的配置结构体指针
	Path       string      // 配置文件路径
	Check      int         // 0表示监听未开启，1表示已开启，2表示停止监听，不会再触发热更新操作
	mu         sync.Mutex
	ch         chan error //传递更新信号或错误的管道
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
		Check:      0,
		mu:         sync.Mutex{},
		ch:         make(chan error),
	}, nil
}

func (l *LocalConfig) GetData() interface{} {
	return l.Data
}

func (l *LocalConfig) LoadData() error {
	l.Viper.SetConfigFile(l.Path)

	// 检查扩展名(支持yaml, json)，如果不是则报错
	ext := filepath.Ext(l.Path)
	switch ext {
	case ".yaml", ".yml":
		l.Viper.SetConfigType("yaml")
	case ".json":
		l.Viper.SetConfigType("json")
	default:
		return errors.New("only .yaml, .yml, or .json are supported")
	}

	// 读取配置文件
	err := loadData(l)
	return err
}

// 传回的通道会在配置文件发生变化时发送信号
func (l *LocalConfig) WatchData() <-chan error {
	// 如果已经监听过就退出，防止开多个协程监听，不将错误传入管道了避免与热更新的错误混淆
	if l.Check == 1 {
		return nil
	}
	l.Check = 1

	l.Viper.OnConfigChange(func(e fsnotify.Event) {
		l.mu.Lock()
		defer l.mu.Unlock()
		// 停止监听就退出
		if l.Check == 2 {
			return
		}
		// 不知道为什么会连着调用两次，所以加个时间限制
		if time.Since(l.Lastupdate) < 1*time.Second {
			return
		}
		err := loadData(l)
		if err != nil {
			// 热更新时出错
			l.ch <- err
			return
		}
		// 发送信号通知配置已更新
		l.ch <- nil
	})
	l.Viper.WatchConfig()

	return l.ch
}

func (l *LocalConfig) Close() error {
	close(l.ch)
	// 停止监听标志
	l.Check = 2
	// 感觉没有错误但得和nacos的一致
	return nil
}

func loadData(l *LocalConfig) error {
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
	Check  int
	ch     chan error
	mu     sync.Mutex
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
		Check:  0,
		ch:     make(chan error),
		mu:     sync.Mutex{},
	}, nil
}

func (p *NacosConfig) GetData() interface{} {
	return p.Data
}

func (p *NacosConfig) LoadData() error {
	content, err := p.Client.GetConfig(vo.ConfigParam{
		DataId: p.DataId,
		Group:  p.Group,
	})
	if err != nil {
		return err
	}
	p.Data = content
	return nil
}

func (p *NacosConfig) WatchData() <-chan error {
	if p.Check == 1 {
		return nil
	}
	p.Check = 1
	err := p.Client.ListenConfig(vo.ConfigParam{
		DataId: p.DataId,
		Group:  p.Group,
		OnChange: func(namespace, group, dataId, data string) {
			p.mu.Lock()
			defer p.mu.Unlock()
			p.Data = data
			p.ch <- nil
		},
	})
	if err != nil {
		p.ch <- err
	}
	return p.ch
}

func (p *NacosConfig) Close() error {
	// nacos提供了专门关闭方法
	cancelParam := vo.ConfigParam{
		DataId: p.DataId,
		Group:  p.Group,
	}
	p.Check = 0
	return p.Client.CancelListenConfig(cancelParam)
}