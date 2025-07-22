package rabbitmq

import (
	"context"
	"fmt"
	"g_mall/pkg/utils/log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ 结构体
type RabbitMQ struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	QueueName    string
	Exchange     string
	Key          string
	MqURL        string
	connNotify   chan *amqp.Error // 连接断开通知
	channNotify  chan *amqp.Error // channel断开通知
	closeFlag    bool             // 是否手动关闭
	mutex        sync.RWMutex     // 读写锁
	retryCount   int              // 重连次数
	maxRetry     int              // 最大重试次数
	reconnecting bool             // 是否正在重连
}

// NewRabbitMQ 创建一个新的操作对象
func NewRabbitMQ(queueName string) *RabbitMQ {
	rabbitMQ := &RabbitMQ{
		QueueName: queueName,
		MqURL:     "amqp://guest:guest@localhost:5672/",
		maxRetry:  5,
		closeFlag: false,
	}

	// 初始化连接
	if err := rabbitMQ.initConnection(); err != nil {
		log.LogrusObj.Errorf("初始化RabbitMQ连接失败: %s", err)
		return nil
	}

	// 启动重连监听
	go rabbitMQ.keepalive()

	return rabbitMQ
}

// initConnection 初始化连接
func (r *RabbitMQ) initConnection() error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	var err error
	r.conn, err = amqp.Dial(r.MqURL)
	if err != nil {
		return fmt.Errorf("创建连接失败: %s", err)
	}

	r.channel, err = r.conn.Channel()
	if err != nil {
		r.conn.Close()
		return fmt.Errorf("创建Channel失败: %s", err)
	}

	// 设置监听器
	r.connNotify = make(chan *amqp.Error, 1)
	r.channNotify = make(chan *amqp.Error, 1)
	r.conn.NotifyClose(r.connNotify)
	r.channel.NotifyClose(r.channNotify)

	// 声明队列
	_, err = r.channel.QueueDeclare(
		r.QueueName, // name
		false,       // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		r.channel.Close()
		r.conn.Close()
		return fmt.Errorf("声明队列失败: %s", err)
	}

	return nil
}

// keepalive 保持连接
func (r *RabbitMQ) keepalive() {
	for {
		if r.closeFlag {
			return
		}

		select {
		case err := <-r.connNotify:
			log.LogrusObj.Errorf("RabbitMQ连接断开: %s", err)
			r.reconnect()
		case err := <-r.channNotify:
			log.LogrusObj.Errorf("RabbitMQ Channel断开: %s", err)
			r.reconnect()
		}
	}
}

// reconnect 重连
func (r *RabbitMQ) reconnect() {
	r.mutex.Lock()
	if r.reconnecting {
		r.mutex.Unlock()
		return
	}
	r.reconnecting = true
	r.mutex.Unlock()

	for i := 0; i < r.maxRetry; i++ {
		if r.closeFlag {
			return
		}

		time.Sleep(time.Second * time.Duration(i+1))
		log.LogrusObj.Infof("尝试第%d次重连", i+1)

		if err := r.initConnection(); err != nil {
			log.LogrusObj.Errorf("第%d次重连失败: %s", i+1, err)
			continue
		}

		log.LogrusObj.Info("重连成功")
		r.mutex.Lock()
		r.reconnecting = false
		r.retryCount = 0
		r.mutex.Unlock()
		return
	}

	log.LogrusObj.Error("重连次数超过最大限制，放弃重连")
	r.mutex.Lock()
	r.reconnecting = false
	r.mutex.Unlock()
}

// PublishSimple 发布简单消息
func (r *RabbitMQ) PublishSimple(ctx context.Context, message []byte) error {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if r.channel == nil {
		return fmt.Errorf("channel未初始化")
	}

	return r.channel.PublishWithContext(ctx,
		"",          // exchange
		r.QueueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message,
		},
	)
}

// ConsumeSimple 消费简单消息
func (r *RabbitMQ) ConsumeSimple(ctx context.Context, handler func(context.Context, []byte) error) {
	msgs, err := r.Consume(ctx)
	if err != nil {
		log.LogrusObj.Errorf("获取消息队列失败: %s", err)
		return
	}

	for msg := range msgs {
		if err := handler(ctx, msg.Body); err != nil {
			log.LogrusObj.Errorf("处理消息失败: %s", err)
		}
	}
}

// Consume 消费消息
func (r *RabbitMQ) Consume(ctx context.Context) (<-chan amqp.Delivery, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	if r.channel == nil {
		return nil, fmt.Errorf("channel未初始化")
	}

	return r.channel.Consume(
		r.QueueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
}

// Destroy 断开channel和connection
func (r *RabbitMQ) Destroy() {
	r.mutex.Lock()
	r.closeFlag = true
	r.mutex.Unlock()

	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}

	log.LogrusObj.Info("RabbitMQ连接已关闭")
}

func ConsumeMessage(ctx context.Context, rabbitMqQueue string) (<-chan amqp.Delivery, error) {
	ch, err := GlobalRabbitMQ.Channel()
	if err != nil {
		fmt.Println("err", err)
	}
	q, _ := ch.QueueDeclare(rabbitMqQueue, false, false, false, false, nil)
	return ch.Consume(q.Name, "", true, false, false, false, nil)
}
