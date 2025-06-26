package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/muxi-Infra/muxi-micro/pkg/config"
)

type DatabaseConfig struct {
	MySQL struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	} `yaml:"mysql"`

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
		DB       string `yaml:"db"`
	} `yaml:"redis"`
}

type Config struct {
	Database DatabaseConfig `yaml:"database"`
}

func Local() {
	// 选择本地就newloaclconfig
	cfg, err := config.NewLocalConfig(&Config{}, "config.yaml")
	if err != nil { //结构体传入时不用指针就报错
		log.Fatal(err)
	}
	// 加载配置
	cfg.LoadData()
	// 检查是否有错误
	if cfg.GetErr() != nil {
		log.Fatalf("加载配置失败: %v", cfg.GetErr())
	}
	// 输出时需要类型断言
	log.Println(cfg.GetData().(*Config).Database.MySQL.Host)

	//ch用于监听热更新信号
	ctx, cancel := context.WithCancel(context.Background())
	ch := cfg.WatchData(ctx)

	go func() {
		for range ch {
			// 用户自定义热更新后的操作
			log.Println(cfg.GetData().(*Config).Database.MySQL.Host)
		}
	}()

	// 热更新出错会直接退出监听并显示错误
	if cfg.GetErr() != nil {
		log.Fatalf("热更新出错: %v", cfg.GetErr())
	}

	time.Sleep(20 * time.Second)
	cancel() // 停止协程
}

func Nacos() {
	// 预先配置
	c := config.NewClientConfig("public", "nacos", "nacos", 5000)
	s := config.NewServerConfig([]string{"localhost"}, []uint64{8848})
	// 选择nacos就NewNacosConfig
	cfg, err := config.NewNacosConfig("key3", "1GRUOP", "value6", c, s)
	if err != nil {
		log.Fatal(err)
	}
    // PUT
	cfg.PutData()
	if cfg.GetErr() != nil {
		log.Fatalf("发布配置失败: %v", cfg.GetErr())
	}

	// GET
	cfg2, _ := config.NewNacosConfig("key3", "1GRUOP", "", c, s)
	cfg2.LoadData()
	if cfg2.GetErr() != nil {
		log.Fatal(cfg2.GetErr())
	}
	log.Println(cfg2.GetData())

	// 监听
	ctx, cancel := context.WithCancel(context.Background())
	ch := cfg2.WatchData(ctx)
	go func() {
		for range ch {
			// 用户自定义热更新后的操作
			log.Println("配置已更新:", cfg2.GetData())
		}
	}()
	// 模拟配置变化
	for i := 0; i < 10; i++ {
		value := fmt.Sprintf("value%d", i)
		cfg, _ = config.NewNacosConfig("key3", "1GRUOP", value, c, s)
		cfg.PutData()
		time.Sleep(5 * time.Second)
	}
	// 监听报错处理
	if cfg2.GetErr() != nil {
		log.Fatal(cfg2.GetErr())
	}

	cancel() // 停止协程
}

func main() {
	Local()
	//Nacos()
}
