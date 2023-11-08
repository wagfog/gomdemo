package main

import (
	"fmt"

	"github.com/micro/go-micro/v2"
)

func main() {
	//服务参数设置
	srv := micro.NewService(
		micro.Name("go.micro.service.user"),
		micro.Version("latest"),
	)
	//初始化服务
	srv.Init()

	// Run service
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
