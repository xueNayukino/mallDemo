package service

import (
	"context"
	"fmt"
	conf "g_mall/config"
	"g_mall/consts"
	"g_mall/pkg/utils/ctl"
	util "g_mall/pkg/utils/log"
	"g_mall/repository/cache"
	"g_mall/repository/db/dao"
	"g_mall/repository/db/model"
	"g_mall/types"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const OrderTimeKey = "OrderTime"

var OrderSrvIns *OrderSrv
var OrderSrvOnce sync.Once

type OrderSrv struct {
}

func GetOrderSrv() *OrderSrv {
	OrderSrvOnce.Do(func() {
		OrderSrvIns = &OrderSrv{}
	})
	return OrderSrvIns
}

func (s *OrderSrv) OrderCreate(ctx context.Context, req *types.OrderCreateReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	order := &model.Order{
		UserID:    u.Id,
		ProductID: req.ProductID,
		BossID:    req.BossID,
		Num:       int(req.Num),
		Money:     float64(req.Money),
		Type:      1,
	}
	addressDao := dao.NewAddressDao(ctx)
	address, err := addressDao.GetAddressByAid(req.AddressID, u.Id)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}

	order.AddressID = address.ID
	number := fmt.Sprintf("%09v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000000))
	productNum := strconv.Itoa(int(req.ProductID))
	userNum := strconv.Itoa(int(u.Id))
	order.OrderNum = number + productNum + userNum

	orderDao := dao.NewOrderDao(ctx)
	err = orderDao.CreateOrder(order)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}

	// 订单号存入Redis中，设置过期时间
	data := redis.Z{
		Score:  float64(time.Now().Unix()) + 15*time.Minute.Seconds(),
		Member: order.OrderNum,
	}
	cache.RedisClient.ZAdd(cache.RedisContext, OrderTimeKey, data)

	return
}

func (s *OrderSrv) OrderList(ctx context.Context, req *types.OrderListReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	orders, total, err := dao.NewOrderDao(ctx).ListOrderByCondition(u.Id, req)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	for i := range orders {
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			orders[i].ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + orders[i].ImgPath
		}
	}

	resp = types.DataListResp{
		Item:  orders,
		Total: total,
	}

	return
}

func (s *OrderSrv) OrderShow(ctx context.Context, req *types.OrderShowReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return nil, err
	}
	order, err := dao.NewOrderDao(ctx).ShowOrderById(req.OrderId, u.Id)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		order.ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + order.ImgPath
	}

	resp = order

	return
}

func (s *OrderSrv) OrderDelete(ctx context.Context, req *types.OrderDeleteReq) (resp interface{}, err error) {
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}
	err = dao.NewOrderDao(ctx).DeleteOrderById(req.OrderId, u.Id)
	if err != nil {
		util.LogrusObj.Error(err)
		return
	}

	return
}
