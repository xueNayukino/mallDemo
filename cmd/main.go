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

//	9     func main() {
//	  10         // 1. 加载配置
//	  11         config.Init()
//	  12
//	  13         // 2. 初始化数据库连接
//	  14         // 从全局配置 Conf 中获取 MySQL 配置
//	  15         dao.InitMySQL(config.Conf.MySQL)
//	  16
//	  17         // 3. 初始化路由
//	  18         r := routes.NewRouter()
//	  19
//	  20         // 4. 启动服务
//	  21         // 从全局配置 Conf 中获取服务端口
//	  22         r.Run(":" + config.Conf.Server.Port)
//	  23     }
func main() {
	loading() // 加载配置
	r := routes.NewRouter()
	_ = r.Run(conf.Config.System.HttpPort)
	fmt.Println("启动配成功...")
}

// loading一些配置
func loading() {
	conf.InitConfig()
	dao.InitMySQL()
	cache.InitCache()
	rabbitmq.InitRabbitMQ()
	rabbitmq.InitSkillProductMQ()
	es.InitEs() // 如果需要接入ELK可以打开这个注释
	//kafka.InitKafka()
	track.InitJaeger()
	util.InitLog() // 如果接入ELK请进入这个func打开注释
	fmt.Println("加载配置完成...")
	go scriptStarting()
}

func scriptStarting() {
	// 启动一些脚本
}
