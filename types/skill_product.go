package types

// SkillProductReq 秒杀请求参数在 skill_goods.go 中已定义

// 新增字段
func (s *SkillProductReq) GetUserId() uint {
	return s.BossId // 在这个场景中，BossId 字段用作用户ID
}

// OrderMessage MQ消息结构
type OrderMessage struct {
	OrderId    string `json:"order_id"`
	UserId     uint   `json:"user_id"`
	ProductId  uint   `json:"product_id"`
	Num        int    `json:"num"`
	CreateTime int64  `json:"create_time"`
}

// OrderStatus 订单状态
const (
	OrderStatusPending = iota // 待处理
	OrderStatusSuccess        // 下单成功
	OrderStatusFailed         // 下单失败
	OrderStatusFinish         // 完成
)
