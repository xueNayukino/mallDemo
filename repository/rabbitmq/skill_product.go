package rabbitmq

import (
	"g_mall/pkg/utils/log"
)

const (
	// 队列名称
	SkillProductQueueName = "skill_product_queue"
	// 死信交换机
	DeadLetterExchange = "skill_product_dlx"
	// 死信队列
	DeadLetterQueue = "skill_product_dlq"
	// 最大重试次数
	MaxRetryCount = 3
)

var (
	SkillProductMQ *RabbitMQ
)

// InitSkillProductMQ 初始化秒杀商品MQ
func InitSkillProductMQ() {
	SkillProductMQ = NewRabbitMQ(SkillProductQueueName)
	if SkillProductMQ == nil {
		log.LogrusObj.Fatal("初始化秒杀商品MQ失败")
	}
}
