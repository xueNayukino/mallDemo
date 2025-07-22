package main

import (
	"fmt"
	conf "g_mall/config"
	util "g_mall/pkg/utils/log"
	"g_mall/pkg/utils/track"
	"g_mall/repository/cache"
	"g_mall/repository/db/dao"
	"g_mall/repository/es"
	"g_mall/repository/rabbitmq"
	"g_mall/routes"
)

func main() {
	loading() // 加载配置
	r := routes.NewRouter()
	_ = r.Run(conf.Config.System.HttpPort)
	fmt.Println("启动配成功...")
}

func loading() {
	conf.InitConfig()
	dao.InitMySQL()
	cache.InitCache()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitSkillProductMQ()
	es.InitEs() 
	track.InitJaeger()
	util.InitLog() 
	fmt.Println("加载配置完成...")
	go scriptStarting()
}

func scriptStarting() {
	// 脚本
}
