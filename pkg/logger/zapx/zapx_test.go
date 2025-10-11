package zapx

import (
	"fmt"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/static"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"testing"
)

func TestNewDefaultZapLogger_AllEnv(t *testing.T) {
	envs := []static.Env{static.EnvDev, static.EnvTest, static.EnvProd}
	for _, env := range envs {
		t.Run(fmt.Sprintf("%v", env), func(t *testing.T) {
			l := NewDefaultZapLogger("./logs/test_default", env)
			if l == nil {
				t.Errorf("NewDefaultZapLogger 返回空")
			}
			logAll(l)
		})
	}
}

func TestNewZapLogger_CustomCore(t *testing.T) {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(devEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.DebugLevel,
	)
	l := NewZapLogger(
		WithZapCore(core),
		WithDefaultZapOptions(),
	)
	if l == nil {
		t.Fatal("自定义 core 创建失败")
	}
	logAll(l)
}

func TestNewZapLogger_WithOutCore(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("core为nil未触发 panic")
		}
	}()

	NewZapLogger()
}

func TestNewZapLogger_CustomOptions(t *testing.T) {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(devEncoderConfig()),
		zapcore.AddSync(os.Stdout),
		zapcore.InfoLevel,
	)
	opts := []zap.Option{zap.AddCaller(), zap.AddCallerSkip(2)}
	l := NewZapLogger(
		WithZapCore(core),
		WithZapOptions(opts...),
	)
	if l == nil {
		t.Fatal("自定义 options 创建失败")
	}
	logAll(l)

}

func TestWithDefaultZapCore_CreatesLogDir(t *testing.T) {
	dir := "./logs/test_create_dir"
	_ = os.RemoveAll(dir)

	opt := WithDefaultZapCore(
		WithCoreEnv(static.EnvTest),
		WithLogDir(dir),
	)
	cfg := &ZapCfg{}
	opt(cfg)

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("日志目录未创建: %v", dir)
	}
}

func TestWithDefaultZapCore_IllegalEnv(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("非法环境未触发 panic")
		}
	}()

	opt := WithDefaultZapCore(WithLogDir("./logs/illegal"), WithCoreEnv(static.Env(99)))
	cfg := &ZapCfg{}
	opt(cfg)
}

func TestWithZapCore_OverridesPrevious(t *testing.T) {
	cfg := &ZapCfg{}
	core1 := zapcore.NewCore(zapcore.NewConsoleEncoder(devEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.InfoLevel)
	core2 := zapcore.NewCore(zapcore.NewJSONEncoder(prodEncoderConfig()), zapcore.AddSync(os.Stdout), zapcore.ErrorLevel)

	WithZapCore(core1)(cfg)
	WithZapCore(core2)(cfg)

	if cfg.core != core2 {
		t.Error("WithZapCore 未正确覆盖 core")
	}
}

func TestWithZapOptions_AppendsOptions(t *testing.T) {
	opt1 := zap.AddCaller()
	opt2 := zap.AddStacktrace(zapcore.DPanicLevel)

	cfg := &ZapCfg{}
	WithZapOptions(opt1, opt2)(cfg)

	if len(cfg.options) != 2 {
		t.Errorf("Option 注入失败，长度应为 2，实际为 %d", len(cfg.options))
	}
}

func TestLogDirClean(t *testing.T) {
	// 测试 logDir clean 是否去除末尾斜杠
	logDir := "./logs/clean-test////"
	opt := WithDefaultZapCore(WithLogDir(logDir), WithCoreEnv(static.EnvTest))
	cfg := &ZapCfg{}
	opt(cfg)

	cleaned := filepath.Clean(logDir)
	if _, err := os.Stat(cleaned); os.IsNotExist(err) {
		t.Errorf("clean 后目录未创建: %s", cleaned)
	}
}

func logAll(l logger.Logger) {

	l.With(logger.Field{
		"string": "string",
		"int":    1,
	})

	l.Info("test",
		logger.Field{
			"string": "string",
			"int":    1,
		},
	)

	l.Debug("test")
	l.Warn("test")
	l.Error("test")
	l.Sync()
}
