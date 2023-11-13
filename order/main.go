package main

import (
	"order/domain/repository"
	"order/domain/service"
	"order/handler"
	order "order/proto"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/registry"
	consul2 "github.com/micro/go-plugins/registry/consul/v2"
	"github.com/micro/go-plugins/wrapper/monitoring/prometheus/v2"
	ratelimit "github.com/micro/go-plugins/wrapper/ratelimiter/uber/v2"
	opentracing2 "github.com/micro/go-plugins/wrapper/trace/opentracing/v2"
	"github.com/opentracing/opentracing-go"
	common "github.com/wagfog/mycommon"
)

const (
	QPS int = 1000
)

func main() {
	consulConfig, err := common.GetConsulConfig("localhost", 8500,
		"/micro/config")
	if err != nil {
		log.Error(err)
	}
	//registry central
	consul := consul2.NewRegistry(func(o *registry.Options) {
		o.Addrs = []string{
			"localhost:8500",
		}
	})

	//jaeger
	t, io, err := common.NewTracer("go.micro.service.order",
		"localhost:6831")
	if err != nil {
		log.Error(err)
	}
	defer io.Close()
	opentracing.SetGlobalTracer(t)

	//init database
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	db.SingularTable(true)

	//tabelInit
	tableinit := repository.NewOrderRepository(db)
	tableinit.InitTable()

	//create instance
	OrderDataService := service.NewOrderDataService(repository.NewOrderRepository(db))

	//expose watch address
	common.PrometheusBoot(9092)

	service := micro.NewService(
		micro.Name("go.micro.service.order"),
		micro.Version("latest"),
		micro.Address(":9085"),
		micro.Registry(consul),
		micro.WrapHandler(opentracing2.NewHandlerWrapper(opentracing.GlobalTracer())),
		micro.WrapHandler(ratelimit.NewHandlerWrapper(QPS)),
		micro.WrapHandler(prometheus.NewHandlerWrapper()),
	)

	service.Init()

	order.RegisterOrderHandler(service.Server(), &handler.Order{OrderDataService: OrderDataService})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
