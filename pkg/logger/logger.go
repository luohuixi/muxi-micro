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

type Field struct {
	Key string
	Val any
}

func Any(key string, val any) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

func Error(err error) Field {
	return Field{
		Key: "error",
		Val: err,
	}
}

func Int64(key string, val int64) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

func Int(key string, val int) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

func String(key string, val string) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

func Int32(key string, val int32) Field {
	return Field{
		Key: key,
		Val: val,
	}
}

// ---------- 环境枚举 ----------
type Env int8

const (
	EnvUnknown Env = iota
	EnvDev         // 开发：彩色多行栈，仅控制台
	EnvTest        // 测试：彩色多行栈到控制台 + JSON 到文件
	EnvProd        // 生产：全 JSON 单行
)
