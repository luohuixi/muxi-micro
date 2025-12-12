package local

import (
	"errors"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/muxi-Infra/muxi-micro/pkg/config"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/logx"
	"github.com/spf13/viper"
)

// OptionStruct 用于存放与泛型 T 无关的配置项，使所有 Option 函数都不需要使用泛型
type OptionStruct struct {
	logger logger.Logger
	size   int
}

type Option func(optionStruct *OptionStruct)

func WithLogger(logger logger.Logger) Option {
	return func(l *OptionStruct) {
		l.logger = logger
	}
}

func WithWatchChanSize(size int) Option {
	return func(l *OptionStruct) {
		l.size = size
	}
}

type LocalConfig[T any] struct {
	viper *viper.Viper
	data  *T
	ch    chan struct{}
	once  sync.Once
	opt   OptionStruct

	sync.RWMutex
}

func LoadLocalConfig[T any](path string, op ...Option) (config.ConfigManager[T], error) {
	v := viper.New()
	v.SetConfigFile(path)
	ext := filepath.Ext(path)
	switch ext {
	case ".yaml", ".yml":
		v.SetConfigType("yaml")
	case ".json":
		v.SetConfigType("json")
	default:
		return nil, errors.New("only .yaml, .yml, or .json are supported")
	}

	var cfg T
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	local := LocalConfig[T]{
		viper: v,
		data:  &cfg,
		opt: OptionStruct{
			logger: logx.NewStdLogger(),
			size:   10,
		},
	}
	for _, o := range op {
		o(&local.opt)
	}

	return &local, nil
}

func (l *LocalConfig[T]) GetData() *T {
	l.RLock()
	defer l.RUnlock()
	return l.data
}

func (l *LocalConfig[T]) WatchData() <-chan struct{} {
	l.once.Do(func() {
		l.ch = make(chan struct{}, l.opt.size)
		l.viper.OnConfigChange(func(e fsnotify.Event) {
			var newData T
			err := l.viper.Unmarshal(&newData)
			if err != nil {
				l.opt.logger.Warn("failed to use viper.Unmarshal while watching", logger.Field{"error": err})
				return
			}
			l.Lock()
			l.data = &newData
			l.Unlock()
			l.opt.logger.Info(fmt.Sprintf("viper find the file change(%s)", e.Op.String()), logger.Field{"file": e.Name})
			l.ch <- struct{}{}
		})
		l.viper.WatchConfig()
	})

	return l.ch
}

func (l *LocalConfig[T]) Close() error {
	l.viper.OnConfigChange(nil)
	if l.ch != nil {
		close(l.ch)
	}
	return nil
}
