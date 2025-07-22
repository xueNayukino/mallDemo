package dao

import (
	"context"
	"g_mall/repository/db/model"
	"g_mall/types"

	"gorm.io/gorm"
)

type ProductDao struct {
	*gorm.DB
}

func NewProductDao(ctx context.Context) *ProductDao {
	return &ProductDao{NewDBClient(ctx)}
}

func NewProductDaoByDB(db *gorm.DB) *ProductDao {
	return &ProductDao{db}
}

// GetProductById 通过 id 获取product
func (dao *ProductDao) GetProductById(id uint) (product *model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where("id=?", id).First(&product).Error

	return
}

// ShowProductById 通过 id 获取product
func (dao *ProductDao) ShowProductById(id uint) (product *model.Product, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where("id=?", id).First(&product).Error

	return
}

// ListProductByCondition 获取商品列表
// SELECT * FROM products WHERE [conditions]
// LIMIT [pageSize] OFFSET [(pageNum-1)*pageSize]
func (dao *ProductDao) ListProductByCondition(condition map[string]interface{}, page types.BasePage) (products []*model.Product, err error) {
	err = dao.DB.Where(condition).
		Offset((page.PageNum - 1) * page.PageSize).
		Limit(page.PageSize).
		Find(&products).Error

	return
}

// CreateProduct 创建商品
// INSERT INTO products (...) VALUES (...)
func (dao *ProductDao) CreateProduct(product *model.Product) error {
	return dao.DB.Model(&model.Product{}).
		Create(&product).Error
}

// CountProductByCondition 根据情况获取商品的数量
// SELECT COUNT(*) FROM products WHERE [conditions]
func (dao *ProductDao) CountProductByCondition(condition map[string]interface{}) (total int64, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where(condition).Count(&total).Error

	return
}

// DeleteProduct 删除商品 pId：商品ID pId：商品ID
// DELETE FROM products WHERE id = ? AND boss_id = ?
func (dao *ProductDao) DeleteProduct(pId, uId uint) error {
	return dao.DB.Model(&model.Product{}).
		Where("id = ? AND boss_id = ?", pId, uId).
		Delete(&model.Product{}).
		Error
}

// UpdateProduct 更新商品
// UPDATE products SET ... WHERE id = ?
func (dao *ProductDao) UpdateProduct(pId uint, product *model.Product) error {
	return dao.DB.Model(&model.Product{}).
		Where("id=?", pId).Updates(&product).Error
}

// SearchProduct 搜索商品
// -- 分页查询
// SELECT * FROM products
// WHERE name LIKE '%keyword%' OR info LIKE '%keyword%'
// LIMIT [pageSize] OFFSET [offset]
//
// -- 总数查询
// SELECT COUNT(*) FROM products
// WHERE name LIKE '%keyword%' OR info LIKE '%keyword%'
func (dao *ProductDao) SearchProduct(info string, page types.BasePage) (products []*model.Product, count int64, err error) {
	err = dao.DB.Model(&model.Product{}).
		Where("name LIKE ? OR info LIKE ?", "%"+info+"%", "%"+info+"%").
		Offset((page.PageNum - 1) * page.PageSize).
		Limit(page.PageSize).
		Find(&products).Error

	if err != nil {
		return
	}

	err = dao.DB.Model(&model.Product{}).
		Where("name LIKE ? OR info LIKE ?", "%"+info+"%", "%"+info+"%").
		Count(&count).
		Error

	return
}
