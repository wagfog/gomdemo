package main

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-plugins/registry/consul/v2"
	"github.com/wagfog/gomdemo/common"
	"github.com/wagfog/gomdemo/domain/repository"
	"github.com/wagfog/gomdemo/domain/service"
	"github.com/wagfog/gomdemo/handler"
	"github.com/wagfog/gomdemo/proto/category"
)

func main() {
	//pei zhi zhong xin
	consulConfig, err := common.GetConsulConfig("127.0.0.1", 8500, "/micro/config")
	if err != nil {
		log.Error(err)
	}
	//registry contre
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs = []string{
			"127.0.0.1:8500",
		}
	})
	//服务参数设置
	srv := micro.NewService(
		micro.Name("go.micro.service.category"),
		micro.Version("latest"),
		//set address and port
		micro.Address("127.0.0.1:8082"),
		micro.Registry(consulRegistry),
	)

	//get mysql option
	mysqlInfo := common.GetMysqlFromConsul(consulConfig, "mysql")
	db, err := gorm.Open("mysql", mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Error(err)
	}
	defer db.Close()
	//禁止复表
	//让 GORM 在操作数据库时不使用复数形式的表名，而是使用单数形式的表名。
	//比如，如果有一个对象叫 "users"，使用这个设置后，它对应的表名会是 "user" 而不是 "users"。
	//这样做可以避免在操作数据库时出现一些表名不一致的问题。
	db.SingularTable(true)
	rp := repository.NewCategoryRepository(db)
	rp.InitTable()
	//初始化服务

	categoryDateService := service.NewCategoryDataService(repository.NewCategoryRepository(db))

	err = category.RegisterCategoryHandler(srv.Server(),
		&handler.Category{CategoryDataService: categoryDateService})
	if err != nil {
		log.Error(err)
	}
	// Run service
	if err := srv.Run(); err != nil {
		fmt.Println(err)
	}
}
