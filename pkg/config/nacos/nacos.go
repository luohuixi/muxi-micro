package nacos

import (
	"encoding/json"
	"sync"

	"github.com/muxi-Infra/muxi-micro/pkg/config"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/logx"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v3"
)

// OptionStruct 用于存放与泛型 T 无关的配置项，使所有 Option 函数都不需要使用泛型
type OptionStruct struct {
	logger logger.Logger
	size   int
	format string
}

type Option func(opt *OptionStruct)

func WithLogger(l logger.Logger) Option {
	return func(opt *OptionStruct) {
		opt.logger = l
	}
}

func WithWatchChanSize(size int) Option {
	return func(opt *OptionStruct) {
		opt.size = size
	}
}

// WithFormat 支持 json/yaml
func WithFormat(format string) Option {
	return func(opt *OptionStruct) {
		opt.format = format
	}
}

type NacosConfig[T any] struct {
	dataId string
	group  string
	val    *T
	client config_client.IConfigClient

	ch   chan struct{}
	once sync.Once
	opt  OptionStruct

	sync.RWMutex
}

type ClientAddr struct {
	Ip   string
	Port uint64
}

func NewClientConfig(namespace, username, password, cacheDir string, timeoutMs uint64) *constant.ClientConfig {
	return &constant.ClientConfig{
		NamespaceId:         namespace,
		TimeoutMs:           timeoutMs,
		Username:            username,
		Password:            password,
		NotLoadCacheAtStart: true,
		CacheDir:            cacheDir,
	}
}

func NewServerConfig(addr []ClientAddr) []constant.ServerConfig {
	serverConfigs := make([]constant.ServerConfig, 0, len(addr))
	for _, a := range addr {
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr: a.Ip,
			Port:   a.Port,
			Scheme: "http",
		})
	}
	return serverConfigs
}

func NewNacosClient(namespace, username, password, cacheDir string, timeoutMs uint64, addr []ClientAddr) (config_client.IConfigClient, error) {
	clientConfig := NewClientConfig(namespace, username, password, cacheDir, timeoutMs)
	serverConfigs := NewServerConfig(addr)

	return clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig:  clientConfig,
		ServerConfigs: serverConfigs,
	})
}

func LoadNacosConfig[T any](client config_client.IConfigClient, group, dataId string, ops ...Option) (config.ConfigManager[T], error) {
	nacos := NacosConfig[T]{
		client: client,
		group:  group,
		dataId: dataId,
		opt: OptionStruct{
			logger: logx.NewStdLogger(),
			size:   10,
			format: "yaml",
		},
	}

	for _, o := range ops {
		o(&nacos.opt)
	}

	content, err := client.GetConfig(vo.ConfigParam{
		DataId: "test",
		Group:  "REPO",
	})
	if err != nil {
		return nil, err
	}

	var cfg T
	if nacos.opt.format == "json" {
		err = json.Unmarshal([]byte(content), &cfg)
	} else {
		err = yaml.Unmarshal([]byte(content), &cfg)
	}
	if err != nil {
		return nil, err
	}

	nacos.val = &cfg
	return &nacos, nil
}

func (nc *NacosConfig[T]) GetData() *T {
	nc.RLock()
	defer nc.RUnlock()
	return nc.val
}

func (nc *NacosConfig[T]) WatchData() <-chan struct{} {
	nc.once.Do(func() {
		nc.ch = make(chan struct{}, nc.opt.size)

		err := nc.client.ListenConfig(vo.ConfigParam{
			DataId: nc.dataId,
			Group:  nc.group,
			OnChange: func(namespace, group, dataId, data string) {
				var cfg T
				var err error

				if nc.opt.format == "json" {
					err = json.Unmarshal([]byte(data), &cfg)
				} else {
					err = yaml.Unmarshal([]byte(data), &cfg)
				}

				if err != nil {
					nc.opt.logger.Error("failed to unmarshal data from nacos", logger.Field{"error": err})
					return
				}

				nc.Lock()
				nc.val = &cfg
				nc.Unlock()
				nc.opt.logger.Info("nacos find the file changed", logger.Field{"dataId": dataId}, logger.Field{"group": group})
				nc.ch <- struct{}{}
			},
		})

		if err != nil {
			nc.opt.logger.Error("Nacos WatchData failed", logger.Field{"error": err})
		}
	})
	return nc.ch
}

func (nc *NacosConfig[T]) Close() error {
	return nc.client.CancelListenConfig(vo.ConfigParam{
		DataId: nc.dataId,
		Group:  nc.group,
	})
}
