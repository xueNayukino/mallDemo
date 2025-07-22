package model

import (
	"g_mall/repository/cache"
	"gorm.io/gorm"
	"strconv"
)

// 商品模型
type Product struct {
	gorm.Model
	Name          string `gorm:"size:255;index"`
	CategoryID    uint   `gorm:"not null"`
	Title         string
	Info          string `gorm:"size:1000"`
	ImgPath       string
	Price         string
	DiscountPrice string
	OnSale        bool `gorm:"default:false"`
	Num           int
	BossID        uint
	BossName      string
	BossAvatar    string
}

// View 获取点击数
func (product *Product) View() uint64 {
	countStr, _ := cache.RedisClient.Get(cache.RedisContext, cache.ProductViewKey(product.ID)).Result()
	count, _ := strconv.ParseUint(countStr, 10, 64)
	return count
}

// AddView 商品游览
func (product *Product) AddView() {
	// 增加视频点击数
	cache.RedisClient.Incr(cache.RedisContext, cache.ProductViewKey(product.ID))
	// 增加排行点击数
	cache.RedisClient.ZIncrBy(cache.RedisContext, cache.RankKey, 1, strconv.Itoa(int(product.ID)))
}

//Incr原子自增1
//2. 更新排行榜：它使用ZINCRBY命令（cache.RedisClient.ZIncrBy(...)）
//对一个名为rank的有序集合（Sorted Set）进行操作。
//* rank是排行榜的Key。
//* 1是这次要增加的分数。
//* strconv.Itoa(int(product.ID))是要更新分数的成员
