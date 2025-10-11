package zapx

import (
	"fmt"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/static"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"os"
	"path/filepath"
)

type ZapLogger struct{ l *zap.Logger }

type ZapOption func(*ZapCfg)

type ZapCfg struct {
	core    zapcore.Core
	options []zap.Option
}

func NewDefaultZapLogger(logDir string, env static.Env) logger.Logger {
	return NewZapLogger(
		WithDefaultZapCore(
			WithLogDir(logDir),
			WithCoreEnv(env),
		),
		WithDefaultZapOptions(),
	)
}

func NewZapLogger(opts ...ZapOption) logger.Logger {
	cfg := &ZapCfg{}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.core == nil {
		log.Panic("缺少 zap-core 核心配置")
	}
	return &ZapLogger{l: zap.New(cfg.core, cfg.options...)}
}

type CoreOption func(*coreCfg)

type coreCfg struct {
	env    static.Env
	logDir string
}

func WithCoreEnv(env static.Env) CoreOption {
	return func(cfg *coreCfg) {
		cfg.env = env
	}
}

func WithLogDir(logDir string) CoreOption {
	return func(cfg *coreCfg) {
		cfg.logDir = logDir
	}
}

func WithDefaultZapCore(opts ...CoreOption) ZapOption {
	return func(cfg *ZapCfg) {
		var corecfg = coreCfg{
			logDir: "./logs",
			env:    static.EnvProd,
		}

		for _, opt := range opts {
			opt(&corecfg)
		}
		// dev 只需要 stdout，不强制创建 logDir
		if corecfg.env != static.EnvDev {
			corecfg.logDir = filepath.Clean(corecfg.logDir)
			if err := os.MkdirAll(corecfg.logDir, 0755); err != nil {
				log.Panicf("无法创建日志目录: %v", err)
			}
		}

		jsonEnc := zapcore.NewJSONEncoder(prodEncoderConfig())
		consoleEnc := zapcore.NewConsoleEncoder(devEncoderConfig())

		switch corecfg.env {
		// ======== DEV：彩色到控制台 ========
		case static.EnvDev:
			cfg.core = zapcore.NewCore(consoleEnc, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
			return

		// ======== TEST：控制台彩色 + 文件 JSON ========
		case static.EnvTest:
			consoleCore := zapcore.NewCore(consoleEnc, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
			fileCore := buildFileCores(jsonEnc, corecfg.logDir, false) // 仅文件
			cfg.core = zapcore.NewTee(append([]zapcore.Core{consoleCore}, fileCore...)...)
			return

		// ======== PROD(默认)：全 JSON 单行 ========
		case static.EnvProd:
			cores := buildFileCores(jsonEnc, corecfg.logDir, true) // stdout+file 共写
			cfg.core = zapcore.NewTee(cores...)

		default:
			log.Panic("非法的环境")
			return
		}
	}
}

// 允许替换 Option
func WithZapOptions(opts ...zap.Option) ZapOption {
	return func(cfg *ZapCfg) { cfg.options = opts }
}

// 默认 Option
func WithDefaultZapOptions() ZapOption {
	return func(cfg *ZapCfg) {
		cfg.options = []zap.Option{
			zap.AddCaller(),
			zap.AddStacktrace(zapcore.WarnLevel),
			zap.AddCallerSkip(1),
		}
	}
}

// 允许自定义 core
func WithZapCore(core zapcore.Core) ZapOption {
	return func(cfg *ZapCfg) { cfg.core = core }
}

// 构造文件相关 core；如果 withStdout=true，则 stdout 也走同 encoder（生产）
func buildFileCores(enc zapcore.Encoder, dir string, withStdout bool) []zapcore.Core {
	var cores []zapcore.Core
	stdout := zapcore.AddSync(os.Stdout)

	var ws zapcore.WriteSyncer
	if withStdout {
		ws = zapcore.NewMultiWriteSyncer(stdout, zapcore.AddSync(newRotateLogger(fmt.Sprintf("%s/app.log", dir))))
	} else {
		ws = zapcore.AddSync(newRotateLogger(fmt.Sprintf("%s/app.log", dir)))
	}
	cores = append(cores, zapcore.NewCore(enc, ws, zapcore.DebugLevel))
	return cores
}

// 生产 JSON encoder
func prodEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:       "@timestamp",
		LevelKey:      "level",
		MessageKey:    "msg",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		EncodeLevel:   zapcore.CapitalLevelEncoder,
		EncodeTime:    zapcore.ISO8601TimeEncoder,
		EncodeCaller:  zapcore.ShortCallerEncoder,
	}
}

// 开发 / 测试 彩色 console encoder
func devEncoderConfig() zapcore.EncoderConfig {
	c := prodEncoderConfig()
	c.EncodeLevel = zapcore.CapitalColorLevelEncoder
	return c
}

// 滚动文件
func newRotateLogger(filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    20,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   true,
	}
}

func (z *ZapLogger) Info(msg string, fields ...logger.Field) { z.l.Info(msg, convert(fields)...) }

func (z *ZapLogger) Error(msg string, fields ...logger.Field) { z.l.Error(msg, convert(fields)...) }

func (z *ZapLogger) Debug(msg string, fields ...logger.Field) { z.l.Debug(msg, convert(fields)...) }

func (z *ZapLogger) Warn(msg string, fields ...logger.Field) { z.l.Warn(msg, convert(fields)...) }

func (z *ZapLogger) Fatal(msg string, fields ...logger.Field) { z.l.Fatal(msg, convert(fields)...) }

func (z *ZapLogger) With(fields ...logger.Field) logger.Logger {
	return &ZapLogger{l: z.l.With(convert(fields)...)}
}

func (z *ZapLogger) Sync() error { return z.l.Sync() }

func convert(fields []logger.Field) []zap.Field {
	var res []zap.Field
	for _, f := range fields {
		for k, v := range f {
			res = append(res, zap.Any(k, v))
		}
	}
	return res
}
