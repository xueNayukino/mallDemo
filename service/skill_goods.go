package service

import (
	"context"
	"encoding/json"
	"fmt"
	"g_mall/pkg/utils/log"
	"g_mall/repository/cache"
	"g_mall/repository/db/dao"
	"g_mall/repository/db/model"
	"g_mall/repository/rabbitmq"
	"g_mall/types"
	"math/rand"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var SkillProductSrvIns *SkillProductSrv
var SkillProductSrvOnce sync.Once

type SkillProductSrv struct {
}

func GetSkillProductSrv() *SkillProductSrv {
	SkillProductSrvOnce.Do(func() {
		SkillProductSrvIns = &SkillProductSrv{}
	})
	return SkillProductSrvIns
}

// InitSkillGoods 初始化商品信息
func (s *SkillProductSrv) InitSkillGoods(ctx context.Context) (resp interface{}, err error) {
	rc := cache.RedisClient
	spList := make([]*model.SkillProduct, 0)
	pipe := rc.Pipeline()

	// 1. 先创建商品列表
	for i := 1; i < 10; i++ {
		sp := &model.SkillProduct{
			ProductId: uint(i),
			BossId:    2,
			Title:     "秒杀商品测试使用",
			Money:     200,
			Num:       10,
		}
		spList = append(spList, sp)

		// 2. 序列化商品信息
		jsonBytes, errx := json.Marshal(sp)
		if errx != nil {
			log.LogrusObj.Errorln("商品序列化失败:", errx)
			return nil, errx
		}

		// 3. 将命令添加到Pipeline
		// 存储商品详情
		pipe.Set(ctx, fmt.Sprintf(cache.SkillProductKey, sp.ProductId), string(jsonBytes), 0)
		// 存储商品库存
		pipe.Set(ctx, fmt.Sprintf(cache.SkillProductStockKey, sp.ProductId), sp.Num, 0)
	}

	// 4. 执行Pipeline
	_, errx := pipe.Exec(ctx)
	if errx != nil {
		log.LogrusObj.Errorln("Redis写入失败:", errx)
		return nil, errx
	}

	// 5. 写入数据库
	if err = dao.NewSkillGoodsDao(ctx).BatchCreate(spList); err != nil {
		log.LogrusObj.Errorln("数据库写入失败:", err)
		return nil, err
	}

	return spList, nil
}

// ListSkillGoods 列表展示
func (s *SkillProductSrv) ListSkillGoods(ctx context.Context) (resp interface{}, err error) {
	rc := cache.RedisClient
	productIds := []uint{1, 2, 3, 4, 5, 6, 7, 8, 9} // 商品ID列表

	skillProducts := make([]*model.SkillProduct, 0)
	pipe := rc.Pipeline()

	// 1. 批量获取所有商品信息和库存
	productCmds := make(map[uint]*redis.StringCmd)
	stockCmds := make(map[uint]*redis.StringCmd)

	for _, pid := range productIds {
		// 获取商品详情
		productCmds[pid] = pipe.Get(ctx, fmt.Sprintf(cache.SkillProductKey, pid))
		// 获取实时库存
		stockCmds[pid] = pipe.Get(ctx, fmt.Sprintf(cache.SkillProductStockKey, pid))
	}

	// 执行pipeline
	_, err = pipe.Exec(ctx)
	if err != nil && err != redis.Nil {
		log.LogrusObj.Errorln("Redis批量读取失败:", err)
		// Redis读取失败，从数据库获取
		return s.getProductsFromDB(ctx)
	}

	// 2. 处理返回结果
	for _, pid := range productIds {
		productJson, err := productCmds[pid].Result()
		if err != nil {
			continue
		}

		var product model.SkillProduct
		if err := json.Unmarshal([]byte(productJson), &product); err != nil {
			log.LogrusObj.Errorln("商品信息解析失败:", err)
			continue
		}

		// 获取实时库存
		if stock, err := stockCmds[pid].Int(); err == nil {
			product.Num = stock
		}

		skillProducts = append(skillProducts, &product)
	}

	// 3. 如果Redis中没有数据，从数据库获取
	if len(skillProducts) == 0 {
		return s.getProductsFromDB(ctx)
	}

	return skillProducts, nil
}

// getProductsFromDB 从数据库获取商品信息并写入缓存
func (s *SkillProductSrv) getProductsFromDB(ctx context.Context) (interface{}, error) {
	rc := cache.RedisClient

	// 1. 从数据库读取
	products, err := dao.NewSkillGoodsDao(ctx).ListSkillGoods()
	if err != nil {
		log.LogrusObj.Errorln("数据库读取失败:", err)
		return nil, err
	}

	// 2. 写入Redis缓存
	pipe := rc.Pipeline()
	for _, product := range products {
		// 序列化商品信息
		jsonBytes, err := json.Marshal(product)
		if err != nil {
			continue
		}

		// 存储商品详情
		pipe.Set(ctx, fmt.Sprintf(cache.SkillProductKey, product.ProductId), string(jsonBytes), 0)
		// 存储商品库存
		pipe.Set(ctx, fmt.Sprintf(cache.SkillProductStockKey, product.ProductId), product.Num, 0)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		log.LogrusObj.Errorln("Redis缓存更新失败:", err)
	}

	return products, nil
}

// GetSkillGoods 详情展示
func (s *SkillProductSrv) GetSkillGoods(ctx context.Context, req *types.GetSkillProductReq) (resp interface{}, err error) {
	rc := cache.RedisClient
	productKey := fmt.Sprintf(cache.SkillProductKey, req.ProductId)

	// 1. 读缓存
	productJson, err := rc.Get(ctx, productKey).Result()
	if err != nil {
		if err == redis.Nil {
			// 缓存未命中，从数据库加载并写入缓存
			product, dbErr := dao.NewSkillGoodsDao(ctx).GetSkillProduct(req.ProductId)
			if dbErr != nil {
				log.LogrusObj.Errorln("数据库查询失败:", dbErr)
				return nil, dbErr
			}

			// 写入缓存
			jsonBytes, _ := json.Marshal(product)
			rc.Set(ctx, productKey, string(jsonBytes), 0)
			rc.Set(ctx, fmt.Sprintf(cache.SkillProductStockKey, product.ProductId), product.Num, 0)

			return product, nil
		} else {
			log.LogrusObj.Errorln("Redis查询失败:", err)
			return nil, err
		}
	}

	// 2. 缓存命中，反序列化返回
	var product model.SkillProduct
	if err = json.Unmarshal([]byte(productJson), &product); err != nil {
		log.LogrusObj.Errorln("商品信息解析失败:", err)
		return nil, err
	}

	return &product, nil
}

// SkillProduct 秒杀商品
func (s *SkillProductSrv) SkillProduct(ctx context.Context, req *types.SkillProductReq) (resp interface{}, err error) {
	rc := cache.RedisClient

	// 1. 生成唯一订单号
	orderId := s.generateOrderId(req.BossId, req.ProductId)

	// 2. 检查是否重复下单
	orderKey := fmt.Sprintf("skill:order:%s", orderId)
	exists, err := rc.Exists(ctx, orderKey).Result()
	if err != nil {
		log.LogrusObj.Errorln("检查订单是否存在失败:", err)
		return nil, fmt.Errorf("系统繁忙，请稍后重试")
	}
	if exists == 1 {
		return nil, fmt.Errorf("订单已存在，请勿重复下单")
	}

	// 3. Redis原子减库存
	stockKey := fmt.Sprintf(cache.SkillProductStockKey, req.ProductId)
	currentStock, err := rc.Decr(ctx, stockKey).Result()
	if err != nil {
		log.LogrusObj.Errorln("Redis扣减库存失败:", err)
		return nil, fmt.Errorf("系统繁忙，请稍后重试")
	}

	// 4. 库存不足，回滚库存
	if currentStock < 0 {
		rc.Incr(ctx, stockKey)
		return nil, fmt.Errorf("商品已售罄")
	}

	// 5. 构建MQ消息
	orderMsg := &types.OrderMessage{
		OrderId:    orderId,
		UserId:     req.BossId,
		ProductId:  req.ProductId,
		Num:        1,
		CreateTime: time.Now().Unix(),
	}

	// 6. 发送消息到MQ
	if err = s.sendToMQ(ctx, orderMsg); err != nil {
		// 发送失败需要回滚库存
		rc.Incr(ctx, stockKey)
		log.LogrusObj.Errorln("发送MQ消息失败:", err)
		return nil, fmt.Errorf("系统繁忙，请稍后重试")
	}

	// 7. 设置订单状态为处理中，过期时间15分钟
	err = rc.Set(ctx, orderKey, types.OrderStatusPending, 15*time.Minute).Err()
	if err != nil {
		log.LogrusObj.Errorln("设置订单状态失败:", err)
		// 这里不需要回滚库存，因为消息已经进入MQ队列
	}

	return orderId, nil
}

// generateOrderId 生成唯一订单号
func (s *SkillProductSrv) generateOrderId(userId, productId uint) string {
	// 简单的订单号生成规则：时间戳+用户ID+商品ID+随机数
	timestamp := time.Now().UnixNano() / 1e6 // 毫秒时间戳
	random := rand.Intn(1000)                // 三位随机数
	return fmt.Sprintf("%d%d%d%03d", timestamp, userId, productId, random)
}

// sendToMQ 发送消息到MQ
func (s *SkillProductSrv) sendToMQ(ctx context.Context, msg *types.OrderMessage) error {
	// 序列化消息
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// 发送到RabbitMQ
	err = rabbitmq.SkillProductMQ.PublishSimple(ctx, msgBytes)
	if err != nil {
		return err
	}

	return nil
}

// SkillProductMQ2MySQL 从mq落库
func SkillProductMQ2MySQL(ctx context.Context, msg []byte) error {
	var orderMsg types.OrderMessage
	if err := json.Unmarshal(msg, &orderMsg); err != nil {
		log.LogrusObj.Errorln("消息解析失败:", err)
		return err
	}

	order := &model.Order{
		UserID:    orderMsg.UserId,
		ProductID: orderMsg.ProductId,
		Num:       orderMsg.Num,
		OrderNum:  orderMsg.OrderId,
		Type:      1, // 假设1代表秒杀订单
	}

	product, err := dao.NewSkillGoodsDao(ctx).GetSkillGood(orderMsg.ProductId)
	if err != nil {
		log.LogrusObj.Errorln("获取商品信息失败:", err)
	} else {
		order.Money = product.Money * float64(order.Num)
	}

	if err = dao.NewOrderDao(ctx).CreateOrder(order); err != nil {
		log.LogrusObj.Errorln("数据库创建订单失败:", err)
		return err
	}

	return nil
}
