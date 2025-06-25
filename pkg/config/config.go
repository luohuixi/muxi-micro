package config

import (
	"context"
	"errors"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/spf13/viper"
)

// yaml和json部分
type ConfigManager struct {
	viper      *viper.Viper
	configMap  map[string]string
	mu         sync.RWMutex
	lastupdate time.Time
}

func NewConfigManager(v *viper.Viper, c map[string]string) *ConfigManager {
	return &ConfigManager{
		viper:      v,
		configMap:  c,
		mu:         sync.RWMutex{},
		lastupdate: time.Now(),
	}
}

// 读取配置文件
func LoadFromLocal(path string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(path)

	// 检查扩展名(支持yaml, json)，如果不是则报错
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		v.SetConfigType("yaml")
	case ".json":
		v.SetConfigType("json")
	default:
		return nil, errors.New("only .yaml, .yml, or .json are supported")
	}

	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}

	return v, nil
}

// 获取配置文件中的某一个的值（如database.mysql.host），值建议存储到哈希方便热更新
func GetConfig(c *viper.Viper, data string) (string, error) {
	// 统一返回string类型，后面按需类型转化
	if !c.IsSet(data) {
		return "", errors.New("key not found: " + data)
	}

	return c.GetString(data), nil
}

// 热更新因为监听的是整个文件，不确定哪个值改了，所有采用哈希的方式，直接修改外部值
func WatchConfig(c *viper.Viper, value map[string]string, ctx context.Context) <-chan struct{} {
	cm := NewConfigManager(c, value)
	ch := make(chan struct{}, 10)

	go func() {
		defer close(ch)
		c.WatchConfig()
		c.OnConfigChange(func(e fsnotify.Event) {
			// 不知道为什么会连着调用两次，所以加个时间限制
			if time.Since(cm.lastupdate) < 1*time.Second {
				return
			}

			cm.reloadAll()
			cm.lastupdate = time.Now()

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

// reloadAll 遍历所有配置项，热更新configMap中的值
func (cm *ConfigManager) reloadAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for _, key := range cm.viper.AllKeys() {
		cm.configMap[key] = cm.viper.GetString(key)
	}
}

// Nacos部分
// 定义查询结构体
type NacosConfig struct {
	DataId  string
	Group   string
	Content string
}

func NewNacosConfig(dataId, group, content string) *NacosConfig {
	return &NacosConfig{
		DataId:  dataId,
		Group:   group,
		Content: content,
	}
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

func (p *NacosConfig) PutToNacos(client config_client.IConfigClient) error {
	_, err := client.PublishConfig(vo.ConfigParam{
		DataId:  p.DataId,
		Group:   p.Group,
		Content: p.Content,
	})
	return err
}

func (p *NacosConfig) GetFromNacos(client config_client.IConfigClient) (string, error) {
	content, err := client.GetConfig(vo.ConfigParam{
		DataId: p.DataId,
		Group:  p.Group,
	})
	if err != nil {
		return "", err
	}
	return content, nil
}

// nacos热更新只监听特定一个键值对，所以不做修改，管道返回新数据，让用户自行修改新数据
func (p *NacosConfig) WatchNacos(client config_client.IConfigClient, ctx context.Context) (<-chan string, <-chan error) {
	ch := make(chan string, 10)
	errch := make(chan error, 1)
	go func() {
		defer func() {
			close(ch)
			close(errch)
		}()
		err := client.ListenConfig(vo.ConfigParam{
			DataId: p.DataId,
			Group:  p.Group,
			OnChange: func(namespace, group, dataId, data string) {
				ch <- data
			},
		})
		if err != nil {
			errch <- err
			return
		}
		<-ctx.Done()
	}()
	return ch, errch
}