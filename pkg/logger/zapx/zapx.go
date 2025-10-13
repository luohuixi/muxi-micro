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

type ZapLogger struct {
	l *zap.Logger
}

type ZapOption func(*ZapCfg)

type ZapCfg struct {
	core    zapcore.Core
	options []zap.Option
	env     static.Env
	logDir  string
}

func NewDefaultZapLogger() logger.Logger {
	return NewZapLogger(
		WithZapCore(NewDefaultZapCore("./logs", static.EnvProd)),
		WithZapOptions(NewDefaultZapOptions()...),
	)
}

func WithCoreEnv(env static.Env) ZapOption {
	return func(cfg *ZapCfg) {
		cfg.env = env
	}
}

// WithZapCore 允许自定义 core, 如果传入了会覆盖Env，logDir的配置，请注意
func WithZapCore(core zapcore.Core) ZapOption {
	return func(cfg *ZapCfg) { cfg.core = core }
}

func WithLogDir(logDir string) ZapOption {
	return func(cfg *ZapCfg) {
		cfg.logDir = logDir
	}
}

// 允许替换 Option
func WithZapOptions(opts ...zap.Option) ZapOption {
	return func(cfg *ZapCfg) { cfg.options = opts }
}

// 默认 core
func NewDefaultZapCore(logDir string, env static.Env) (core zapcore.Core) {

	// dev 只需要 stdout，不强制创建 logDir
	if env != static.EnvDev {
		logDir = filepath.Clean(logDir)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Panicf("无法创建日志目录: %v", err)
		}
	}

	jsonEnc := zapcore.NewJSONEncoder(prodEncoderConfig())
	consoleEnc := zapcore.NewConsoleEncoder(devEncoderConfig())

	switch env {
	// ======== DEV：彩色到控制台 ========
	case static.EnvDev:
		core = zapcore.NewCore(consoleEnc, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
		return

	// ======== TEST：控制台彩色 + 文件 JSON ========
	case static.EnvTest:
		consoleCore := zapcore.NewCore(consoleEnc, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)
		fileCore := buildFileCores(jsonEnc, logDir, false) // 仅文件
		core = zapcore.NewTee(append([]zapcore.Core{consoleCore}, fileCore...)...)
		return

	// ======== PROD(默认)：全 JSON 单行 ========
	case static.EnvProd:
		cores := buildFileCores(jsonEnc, logDir, true) // stdout+file 共写
		core = zapcore.NewTee(cores...)

	default:
		log.Panic("非法的环境")
		return
	}

	return core
}

// 默认 Option
func NewDefaultZapOptions() []zap.Option {
	return []zap.Option{
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.WarnLevel),
		zap.AddCallerSkip(1),
	}
}

func NewZapLogger(opts ...ZapOption) logger.Logger {
	cfg := &ZapCfg{
		logDir: "./logs",
		env:    static.EnvProd,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// 如果用户没有传入 core，则使用默认的
	if cfg.core == nil {
		cfg.core = NewDefaultZapCore(cfg.logDir, cfg.env)
	}

	return &ZapLogger{l: zap.New(cfg.core, cfg.options...)}
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
