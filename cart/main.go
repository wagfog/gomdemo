package main

import (
	"github.com/wagfog/micro_cart/common"
	"github.com/wagfog/micro_cart/domain/repository"
	"github.com/wagfog/micro_cart/handler"
	cart "github.com/wagfog/micro_cart/proto"

	service2 "github.com/wagfog/micro_cart/domain/service"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
)

var QPS = 100

func main() {
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}

	//registry
	consul := consul2.NewRegistry(func(o *registry.Options) {
		o.Addrs = []string{
			"127.0.0.1:8500",
		}
	})

	t, io, err := common.NewTracer("go.micro.service.cart", "localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//数据库连接
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	//创建数据库连接
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	//禁止副表
	db.SingularTable(true)

	//第一次初始化
	err = repository.NewCartRepository(db).InitTable()
	if err != nil {
		log.Error(err)
	}

	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.cart"),
		micro.Version("latest"),
		//暴露的服务地址
		micro.Address("0.0.0.0:8087"),
		//注册中心
		micro.Registry(consul),
		//链路追踪
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		//添加限流
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
	)

	// Initialise service
	service.Init()

	cartDataService := service2.NewCartDataService(repository.NewCartRepository(db))

	// Register Handler
	cart.RegisterCartHandler(service.Server(), &handler.Cart{CartDataService: cartDataService})

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
