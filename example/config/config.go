package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/muxi-Infra/muxi-micro/pkg/config"
)

func Local() {
	c, err := config.LoadFromLocal("config.yaml")

	if err != nil {
		log.Fatal(err)
	}

	m := make(map[string]string, 10)
	key := "database.mysql.host"

	//统一将值存到哈希
	value, err := config.GetConfig(c, key)
	m[key] = value

	if err != nil {
		log.Fatal(err)
	}

	log.Println(m[key])

	//ch用于监听热更新信号，stop需在程序结束前手动关闭
	ctx, cancel := context.WithCancel(context.Background())
	ch := config.WatchConfig(c, m, ctx)

	go func() {
		for range ch {
			// 用户自定义热更新后的操作
			log.Println(m[key])
		}
	}()

	time.Sleep(20 * time.Second)
	cancel() // 停止协程
}

func Nacos() {
	c := config.NewClientConfig("public", "nacos", "nacos", 5000)
	s := config.NewServerConfig([]string{"localhost"}, []uint64{8848})
	client, err := config.NewNacos(c, s)
	if err != nil {
		log.Fatalf("创建Nacos客户端失败: %v", err)
	}

	// PUT
	check := config.NewNacosConfig("key3", "1GRUOP", "value6")
	err = check.PutToNacos(client)
	if err != nil {
		log.Fatalf("发布配置失败: %v", err)
	}

	// GET
	check = config.NewNacosConfig("key3", "1GRUOP", "")
	content, err := check.GetFromNacos(client)
	if err != nil {
		log.Fatalf("获取配置失败: %v", err)
	}
	log.Println("获取到的配置内容:", content)

	// 监听
	ctx, cancel := context.WithCancel(context.Background())
	check = config.NewNacosConfig("key3", "1GRUOP", "")
	ch, errch := check.WatchNacos(client, ctx)
	go func() {
		for data := range ch {
			// 用户自定义热更新后的操作
			log.Println("配置已更新:", data)
		}
	}()
	go func() {
		for err := range errch {
			// 用户自定义错误处理
			log.Println("监听配置变化时发生错误:", err)
		}
	}()

	// 模拟配置变化
	for i := 0; i < 10; i++ {
		value := fmt.Sprintf("value%d", i)
		check := config.NewNacosConfig("key3", "1GRUOP", value)
		_ = check.PutToNacos(client)
		time.Sleep(5 * time.Second) 
	}

	cancel() // 停止协程
}

func main(){
	// yaml和json示例
	Local()

	// Nacos示例
	//Nacos()
}