package main

import (
	"fmt"

	"github.com/micro/go-micro/v2"
)

func main() {
	// Create service
	srv := micro.NewService(
		micro.Name("go.micro.service.category"),
		micro.Version("latest"),
		//set address and port
	)

	srv.Init()

	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
