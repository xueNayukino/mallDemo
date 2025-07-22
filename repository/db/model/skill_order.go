package model

import "time"

// SkillOrder 秒杀订单模型
type SkillOrder struct {
	ID        uint   `gorm:"primarykey"`
	OrderId   string `gorm:"type:varchar(50);uniqueIndex;not null"`
	UserId    uint   `gorm:"not null"`
	ProductId uint   `gorm:"not null"`
	Num       int    `gorm:"not null"`
	Status    int    `gorm:"not null;default:1"` // 1:创建成功 2:已支付 3:已取消
	CreatedAt time.Time
	UpdatedAt time.Time
}
