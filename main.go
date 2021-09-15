package main

import (
	"category/common"
	"category/domain/repository"
	service2 "category/domain/service"
	"category/handler"
	category "category/proto/category"
	"github.com/asim/go-micro/plugins/registry/consul/v3"
	"github.com/asim/go-micro/v3"
	log "github.com/asim/go-micro/v3/logger"
	"github.com/asim/go-micro/v3/registry"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	//配置中心
	consulConfig,err := common.GetConsulConfig("127.0.0.1",8500,"/micro/config")
	if err!=nil{
		log.Error(err)
	}
	//注册中心
	consulRegistry := consul.NewRegistry(func(options *registry.Options) {
		options.Addrs=[]string{
			"127.0.0.1:8500",
		}
	})
	//New Service
	service := micro.NewService(
		micro.Name("go.micro.service.category"),
		micro.Version("latest"),
		//这里设置地址和需要暴露的端口
		micro.Address("127.0.0.1:8082"),
		//添加concul作为注册中心
		micro.Registry(consulRegistry),
		)

	//获取mysql配置,路径中不带前缀  //数数据库中间件
	mysqlInfo := common.GetMysqlFromConsul(consulConfig,"mysql")
	//链接数据库
	db,err := gorm.Open("mysql",mysqlInfo.User+":"+mysqlInfo.Pwd+"@/"+mysqlInfo.Database+"?charset=utf8&parseTime=True&loc=Local")
	if err !=nil {
		log.Error(err)
	}
	defer db.Close()
	//禁止复表
	db.SingularTable(true)

	//创建表 只执行一次
	//rp := repository.NewCategoryRepository(db)
	//rp.InitTable()

	//Initialise service
	service.Init()
	//实例化
	categoryDataService := service2.NewCategoryDataService(repository.NewCategoryRepository(db))
	//Register Handler
	err = category.RegisterCategoryHandler(service.Server(),&handler.Category{CategoryDataService:categoryDataService})
	//category.RegisterCategoryHandler(service.Server(),new(handler.Category))
	if err!=nil{
		log.Error(err)
	}
	//Run service
	if err:=service.Run();err!=nil{
		log.Fatal(err)
	}
}
