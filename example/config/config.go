package main

import (
	"log"

	"github.com/muxi-Infra/muxi-micro/pkg/config/local"
	"github.com/muxi-Infra/muxi-micro/pkg/config/nacos"
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
	cfg, err := local.LoadLocalConfig[Config]("./example.yaml", local.WithWatchChanSize(10))
	if err != nil {
		log.Fatal(err)
	}
	defer cfg.Close()

	data := cfg.GetData()
	log.Printf("MySQL Host: %s", data.Database.MySQL.Host)
	log.Printf("Redis DB: %s", data.Database.Redis.DB)

	ch := cfg.WatchData()
	for range ch {
		updatedData := cfg.GetData()
		log.Printf("MySQL Host: %s", updatedData.Database.MySQL.Host)
	}
}

func Nacos() {
	clientAddr := []nacos.ClientAddr{
		{Ip: "localhost", Port: 8848},
	}
	client, err := nacos.NewNacosClient(
		"new",
		"nacos",
		"nacos",
		"",
		1000,
		clientAddr,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer client.CloseClient()

	cfg, err := nacos.LoadNacosConfig[Config](client, "test", "test")
	if err != nil {
		log.Fatal(err)
	}
	defer cfg.Close()

	data := cfg.GetData()
	log.Printf("MySQL Host: %s", data.Database.MySQL.Host)
	log.Printf("Redis DB: %s", data.Database.Redis.DB)

	ch := cfg.WatchData()
	for range ch {
		updatedData := cfg.GetData()
		log.Printf("MySQL Host: %s", updatedData.Database.MySQL.Host)
	}
}

func main() {
	Local()
	//Nacos()
}
