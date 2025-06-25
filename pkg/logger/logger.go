package logger

type Logger interface {
	Info(msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	With(fields ...Field) Logger
	Sync() error
}

type Field interface{}

// ---------- 环境枚举 ----------
type Env int8

const (
	EnvUnknown Env = iota
	EnvDev         // 开发：彩色多行栈，仅控制台
	EnvTest        // 测试：彩色多行栈到控制台 + JSON 到文件
	EnvProd        // 生产：全 JSON 单行
)
