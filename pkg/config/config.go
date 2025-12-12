package config

type ConfigManager[T any] interface {
	GetData() *T                // 获取配置数据
	WatchData() <-chan struct{} // 热更新, 只负责更新数据, 不负责更新后的操作
	Close() error               // 关闭
}
