package dao

import (
	"context"
	"g_mall/consts"
	"g_mall/repository/db/model"

	"gorm.io/gorm"
)

type SkillGoodsDao struct {
	*gorm.DB
}

func NewSkillGoodsDao(ctx context.Context) *SkillGoodsDao {
	return &SkillGoodsDao{NewDBClient(ctx)}
}

func (dao *SkillGoodsDao) Create(in *model.SkillProduct) error {
	return dao.Model(&model.SkillProduct{}).Create(&in).Error
}

func (dao *SkillGoodsDao) BatchCreate(in []*model.SkillProduct) error {
	return dao.Model(&model.SkillProduct{}).
		CreateInBatches(&in, consts.ProductBatchCreate).Error
}

func (dao *SkillGoodsDao) CreateByList(in []*model.SkillProduct) error {
	return dao.Model(&model.SkillProduct{}).Create(&in).Error
}

func (dao *SkillGoodsDao) GetSkillProduct(productId uint) (*model.SkillProduct, error) {
	var product model.SkillProduct
	err := dao.Model(&model.SkillProduct{}).Where("product_id = ?", productId).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// CheckOrderExists 检查订单是否存在
func (dao *SkillGoodsDao) CheckOrderExists(orderId string) (bool, error) {
	var count int64
	err := dao.Model(&model.SkillOrder{}).Where("order_id = ?", orderId).Count(&count).Error
	return count > 0, err
}

// CreateOrder 创建秒杀订单
func (dao *SkillGoodsDao) CreateOrder(orderId string, userId uint, productId uint, num int) error {
	// 创建秒杀订单
	order := &model.SkillOrder{
		OrderId:   orderId,
		UserId:    userId,
		ProductId: productId,
		Num:       num,
		Status:    1, // 1: 创建成功
	}

	// 开启事务
	tx := dao.Begin()

	// 1. 创建订单
	if err := tx.Create(order).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 2. 扣减库存
	if err := tx.Model(&model.SkillProduct{}).
		Where("product_id = ? AND num >= ?", productId, num).
		UpdateColumn("num", tx.Raw("num - ?", num)).Error; err != nil {
		tx.Rollback()
		return err
	}

	// 提交事务
	return tx.Commit().Error
}

func (dao *SkillGoodsDao) ListSkillGoods() (resp []*model.SkillProduct, err error) {
	err = dao.Model(&model.SkillProduct{}).
		Where("num > 0").Find(&resp).Error

	return
}

func (dao *SkillGoodsDao) GetSkillGood(id uint) (resp *model.SkillProduct, err error) {
	err = dao.Model(&model.SkillProduct{}).
		Where("product_id = ?", id).First(&resp).Error

	return
}
