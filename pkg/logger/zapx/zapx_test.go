package zapx

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestNewDefaultZapLogger_AllEnv(t *testing.T) {
	envs := []logger.Env{logger.EnvDev, logger.EnvTest, logger.EnvProd}
	for _, env := range envs {
		t.Run(fmt.Sprintf("%v", env), func(t *testing.T) {
			l := NewDefaultZapLogger("./logs/test_default", env)
			if l == nil {
				t.Errorf("NewDefaultZapLogger 返回空")
			}
			logAll(l, t)
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
	logAll(l, t)
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
	logAll(l, t)

}

func TestWithDefaultZapCore_CreatesLogDir(t *testing.T) {
	dir := "./logs/test_create_dir"
	_ = os.RemoveAll(dir)

	opt := WithDefaultZapCore(
		WithCoreEnv(logger.EnvTest),
		WithCoreSplit(false),
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

	opt := WithDefaultZapCore(WithCoreSplit(true), WithLogDir("./logs/illegal"), WithCoreEnv(logger.Env(99)))
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
	opt := WithDefaultZapCore(WithCoreSplit(true), WithLogDir(logDir), WithCoreEnv(logger.EnvTest))
	cfg := &ZapCfg{}
	opt(cfg)

	cleaned := filepath.Clean(logDir)
	if _, err := os.Stat(cleaned); os.IsNotExist(err) {
		t.Errorf("clean 后目录未创建: %s", cleaned)
	}
}

func logAll(l logger.Logger, t *testing.T) {

	l.With(logger.String("all", "everyLog"))
	l.Info("test", logger.Int("1", 1))
	l.Debug("test", logger.Int("1", 1))
	l.Warn("test", logger.Int("1", 1))
	l.Error("test", logger.Int("1", 1))
	l.Sync()
}
