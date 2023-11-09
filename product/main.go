package main

import (
	"fmt"
	"product/common"
	"product/domain/repository"
	service2 "product/domain/service"
	"product/handler"
	"product/proto/product"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	registry "github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-plugins/registry/consul/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

func main() {
	//option central
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	//registry central
	consul := consul.NewRegistry(func(o *registry.Options) {
		o.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	//链路追踪
	t, io, err := common.NewTracer("go.micro.service.product", "localhost:6831")
	if err != nil {
		log.Fatal(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//database set
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")

	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=true&loc=Local")

	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	db.SingularTable(true)
	//init
	// repository.NewProductRepository(db).InitTable()
	// Create service
	productDataService := service2.NewProductDataService(repository.NewProductRepository(db))

	//set service

	srv := micro.NewService(
		micro.Name("go.micro.service.product"),
		micro.Version("latest"),
		micro.Address("127.0.0.1:8082"),
		//set registry
		micro.Registry(consul),
		//绑定链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
	)

	srv.Init()
	product.RegisterProductHandler(srv.Server(), &handler.Product{ProductDataService: productDataService})
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
