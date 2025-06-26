package main

import (
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
	err = cfg.LoadData()
	// 检查是否有错误
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	// 输出时需要类型断言
	log.Println(cfg.GetData().(*Config).Database.MySQL.Host)

	//ch用于监听热更新信号
	ch := cfg.WatchData()

	go func() {
		for err := range ch {
			// 用户自定义热更新后的操作
			if err == nil {
				log.Println(cfg.GetData().(*Config).Database.MySQL.Host)
			} else {
				log.Println("err: ", err)
			}
		}
	}()

	time.Sleep(20 * time.Second)
	// 其实local的关闭不会报错，nacos才会
	err = cfg.Close()
	if err != nil {
		log.Println(err)
	}
}

func Nacos() {
	// 预先配置
	c := config.NewClientConfig("public", "nacos", "nacos", 5000)
	s := config.NewServerConfig([]string{"localhost"}, []uint64{8848})
	// GET
	cfg2, err := config.NewNacosConfig("key3", "1GRUOP", "", c, s)
	if err != nil {
		log.Fatal(err)
	}
	err = cfg2.LoadData()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(cfg2.GetData())

	// 监听
	ch := cfg2.WatchData()
	go func() {
		for err := range ch {
			// 用户自定义热更新后的操作
			if err != nil {
				log.Println(err)
			} else {
				log.Println("配置已更新:", cfg2.GetData())
			}
		}
	}()

	time.Sleep(10 * time.Second)
	err = cfg2.Close()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	Local()
	//Nacos()
}
