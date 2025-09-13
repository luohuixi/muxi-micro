package infrastructure

import (
	"errors"
	"github.com/google/wire"
	"helloworld/configs"
	"helloworld/internal/repository"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/zapx"
	"sync"
	"time"
)

var ProviderSet = wire.NewSet(
	ProvideGrpcInstance,
	ProvideDB,
	ProvideCache,
)

type Grpc struct {
	logger  logger.Logger
	Addr    string
	Timeout time.Duration
	Net     string
}

var (
	Err      error
	Once     sync.Once
	instance *Grpc
)

var MissingAddrErr = errors.New("you should add a address(like 0.0.0.0:50051) on the configs")

func ProvideGrpcInstance() (*Grpc, error) {
	Once.Do(func() {
		cfg, err := configs.LoadConfig()
		if err != nil {
			Err = err
			return
		}
		if cfg.Server.Grpc.Addr == "" {
			Err = MissingAddrErr
			return
		}
		if cfg.Server.Grpc.Timeout == 0 {
			cfg.Server.Grpc.Timeout = 10 * time.Second
		}
		if cfg.Server.Grpc.Network == "" {
			cfg.Server.Grpc.Network = "tcp"
		}
		instance = &Grpc{
			logger:  zapx.NewDefaultZapLogger(cfg.Log.Dir, logger.EnvTest),
			Addr:    cfg.Server.Grpc.Addr,
			Timeout: cfg.Server.Grpc.Timeout,
			Net:     cfg.Server.Grpc.Network,
		}
	})
	return instance, Err
}

func ProvideDB() (string, error) {
	cfg, err := configs.LoadConfig()
	if err != nil {
		return "", err
	}
	return cfg.Data.MySQL.Dsn, nil
}

func ProvideCache(g *Grpc) (*repository.CacheStruct, error) {
	cfg, err := configs.LoadConfig()
	if err != nil {
		return nil, err
	}
	data := cfg.Data.Redis
	return &repository.CacheStruct{
		RedisAddr:     data.Addr,
		RedisPassword: data.Password,
		Number:        data.Num,
		TtlForCache:   data.Read,
		TtlForSet:     data.Write,
		Log:           g.logger,
	}, nil
}
