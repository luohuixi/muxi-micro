package main

import (
	"context"
	"fmt"
	"github.com/muxi-Infra/muxi-micro/tool/curd/template"
	"time"
)

func main() {
	ttl := 10 * time.Second
	instance, err := template.NewUserModels(
		"root:2388287244@tcp(112.126.68.22:3306)/library?parseTime=true&charset=utf8mb4&loc=Local",
		"112.126.68.22:6379",
		"lhx2388287244",
		0,
		-1,
		ttl,
	)
	if err != nil {
		fmt.Println(err)
	}
	user := template.User{
		Password: "666",
		Username: "test777",
		Mobile:   "1145141919810",
	}
	err = instance.Create(context.Background(), &user)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(user)

	select {}
}
