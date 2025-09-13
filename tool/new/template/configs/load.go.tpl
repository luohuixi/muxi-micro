package configs

import (
	"github.com/muxi-Infra/muxi-micro/pkg/config"
	"sync"
	"time"
)

var (
	configOnce     sync.Once
	configInstance *Config
	configErr      error
)

type Config struct {
	Server struct {
		Grpc struct {
			Addr    string        `yaml:"addr"`
			Network string        `yaml:"network"`
			Timeout time.Duration `yaml:"timeout"`
		} `yaml:"grpc"`
	} `yaml:"server"`
	Data struct {
		MySQL struct {
			Dsn string `yaml:"dsn"`
		} `yaml:"mysql"`
		Redis struct {
			Addr     string        `yaml:"addr"`
			Password string        `yaml:"password"`
			Num      int           `yaml:"num"`
			Read     time.Duration `yaml:"read"`
			Write    time.Duration `yaml:"write"`
		} `yaml:"redis"`
	} `yaml:"data"`
	Log struct {
		Dir string `yaml:"dir"`
	} `yaml:"log"`
}

func LoadConfig() (*Config, error) {
	configOnce.Do(func() {
		cfg, err := config.NewLocalConfig(&Config{}, "../../configs/config.yaml")
		if err != nil {
			configErr = err
			return
		}
		err = cfg.LoadData()
		if err != nil {
			configErr = err
			return
		}
		configInstance = cfg.GetData().(*Config)
	})
	return configInstance, configErr
}
