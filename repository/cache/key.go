package cache

import (
	"fmt"
	"strconv"
)

const (
	// RankKey 每日排名
	RankKey              = "rank"
	SkillProductKey      = "skill:product:%d"       // 商品详情
	SkillProductStockKey = "skill:product:stock:%d" // 商品库存
	SkillProductListKey  = "skill:product_list"
	SkillProductUserKey  = "skill:user:%s"
)

// // 生成商品浏览量的Redis键
func ProductViewKey(id uint) string {
	return fmt.Sprintf("view:product:%s", strconv.Itoa(int(id)))
}
