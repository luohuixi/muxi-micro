package main

import (
	"context"
	"fmt"
	"github.com/muxi-Infra/muxi-micro/pkg/logger"
	"github.com/muxi-Infra/muxi-micro/pkg/logger/zapx"
	"github.com/muxi-Infra/muxi-micro/tool/curd/example"
	"log"
	"time"
)

func main() {
	DBdsn := "your db dsn"
	redisAddr := "your redis addr"
	redisPassword := "your redis password"
	redisDB := 0
	// 缓存持续时间
	ttlForCache := 5 * time.Second
	// 异步设置缓存时间
	ttlForSet := 5 * time.Second
	// 日志记录 redis 错误
	l := zapx.NewDefaultZapLogger("./logs", logger.EnvTest)
	instance, err := example.NewUserModels(
		DBdsn,
		redisAddr,
		redisPassword,
		redisDB,
		ttlForCache,
		ttlForSet,
		l,
	)
	if err != nil {
		log.Fatal(err)
	}
	//user := example.User{
	//	Id:       7,
	//	Username: "example3",
	//	Password: "123456",
	//	Mobile:   "example111",
	//}
	//_ = instance.Create(context.Background(), &user)
	//_ = instance.Update(context.Background(), &user)
	//_ = instance.Delete(context.Background(), 7)
	//value, err := instance.FindOne(context.Background(), 3)
	value, err := instance.FindByMobile(context.Background(), "123456")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)
}
