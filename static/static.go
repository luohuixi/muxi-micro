package static

// ---------- 环境枚举 ----------
type Env int8

const (
	EnvUnknown Env = iota
	EnvDev
	EnvTest
	EnvProd
)
