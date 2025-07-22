package service

import (
	"context"
	conf "g_mall/config"
	"g_mall/consts"
	"g_mall/pkg/utils/ctl"
	"g_mall/pkg/utils/log"
	util "g_mall/pkg/utils/upload"
	"g_mall/repository/db/dao"
	"g_mall/repository/db/model"
	"g_mall/types"
	"mime/multipart"
	"strconv"
	"sync"
)

var ProductSrvIns *ProductSrv
var ProductSrvOnce sync.Once

type ProductSrv struct {
}

func GetProductSrv() *ProductSrv {
	ProductSrvOnce.Do(func() {
		ProductSrvIns = &ProductSrv{}
	})
	return ProductSrvIns
}

// ProductShow 商品
func (s *ProductSrv) ProductShow(ctx context.Context, req *types.ProductShowReq) (resp interface{}, err error) {
	p, err := dao.NewProductDao(ctx).ShowProductById(req.ID) //返回product *model.Product
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	pResp := &types.ProductResp{
		ID:            p.ID,
		Name:          p.Name,
		CategoryID:    p.CategoryID,
		Title:         p.Title,
		Info:          p.Info,
		ImgPath:       p.ImgPath,
		Price:         p.Price,
		DiscountPrice: p.DiscountPrice,
		View:          p.View(), //这个是model层的一个方法，获得点击数
		CreatedAt:     p.CreatedAt.Unix(),
		Num:           p.Num,
		OnSale:        p.OnSale,
		BossID:        p.BossID,
		BossName:      p.BossName,
		BossAvatar:    p.BossAvatar,
	}
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		//// 3. 如果是本地存储模式，拼接完整图片URL
		pResp.BossAvatar = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.AvatarPath + pResp.BossAvatar
		pResp.ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + pResp.ImgPath
	}

	resp = pResp

	return
}

// 创建商品，files：上传的文件列表（multipart.FileHeader 数组）
func (s *ProductSrv) ProductCreate(ctx context.Context, files []*multipart.FileHeader, req *types.ProductCreateReq) (resp interface{}, err error) {
	/// 1. 从上下文中获取用户信息
	u, err := ctl.GetUserInfo(ctx)
	if err != nil {
		log.LogrusObj.Error(err)
		return nil, err
	}
	uId := u.Id
	// 2. 获取商家信息
	boss, _ := dao.NewUserDao(ctx).GetUserById(uId) //也就是用户信息
	// 以第一张作为封面图
	tmp, _ := files[0].Open()
	var path string
	if conf.Config.System.UploadModel == consts.UploadModelLocal {
		// 上传到本地
		path, err = util.ProductUploadToLocalStatic(tmp, uId, req.Name)
	} else {
		// 上传到七牛云
		path, err = util.UploadToQiNiu(tmp, files[0].Size)
	}
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	// 4. 创建商品记录
	product := &model.Product{
		Name:          req.Name,
		CategoryID:    req.CategoryID,
		Title:         req.Title,
		Info:          req.Info,
		ImgPath:       path, // 封面图路径
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		Num:           req.Num,
		OnSale:        true, // 默认上架
		BossID:        uId,
		BossName:      boss.UserName,
		BossAvatar:    boss.Avatar,
	}
	// 5. 保存到数据库
	productDao := dao.NewProductDao(ctx)
	err = productDao.CreateProduct(product)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}

	// 6. 使用WaitGroup并发处理其他图片
	wg := new(sync.WaitGroup)
	wg.Add(len(files)) // 设置等待的协程数

	// 并发上传其他图片并保存到数据库
	for index, file := range files {
		num := strconv.Itoa(index)
		tmp, _ = file.Open()
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			path, err = util.ProductUploadToLocalStatic(tmp, uId, req.Name+num)
		} else {
			path, err = util.UploadToQiNiu(tmp, file.Size)
		}
		if err != nil {
			log.LogrusObj.Error(err)
			return
		}
		// 创建商品图片记录
		productImg := &model.ProductImg{
			ProductID: product.ID,
			ImgPath:   path,
		}
		// 保存图片记录到数据库
		err = dao.NewProductImgDaoByDB(productDao.DB).CreateProductImg(productImg)
		if err != nil {
			log.LogrusObj.Error(err)
			return
		}
		wg.Done() // 标记一个图片处理完成
	}

	wg.Wait()

	return
}

// 查询商品列表
func (s *ProductSrv) ProductList(ctx context.Context, req *types.ProductListReq) (resp interface{}, err error) {
	var total int64
	condition := make(map[string]interface{}) //查询条件
	if req.CategoryID != 0 {
		condition["category_id"] = req.CategoryID //添加分类条件
	}

	productDao := dao.NewProductDao(ctx)
	products, _ := productDao.ListProductByCondition(condition, req.BasePage) //商品列表
	total, err = productDao.CountProductByCondition(condition)                //商品数
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	//
	pRespList := make([]*types.ProductResp, 0)
	for _, p := range products {
		pResp := &types.ProductResp{
			ID:            p.ID,
			Name:          p.Name,
			CategoryID:    p.CategoryID,
			Title:         p.Title,
			Info:          p.Info,
			ImgPath:       p.ImgPath,
			Price:         p.Price,
			DiscountPrice: p.DiscountPrice,
			View:          p.View(),
			CreatedAt:     p.CreatedAt.Unix(),
			Num:           p.Num,
			OnSale:        p.OnSale,
			BossID:        p.BossID,
			BossName:      p.BossName,
			BossAvatar:    p.BossAvatar,
		}
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			pResp.BossAvatar = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.AvatarPath + pResp.BossAvatar
			pResp.ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + pResp.ImgPath
		}
		pRespList = append(pRespList, pResp) //每个商品 添加到响应列表
	}

	resp = &types.DataListResp{
		Item:  pRespList,
		Total: total,
	}

	return
}

// ProductDelete 删除商品，指定
func (s *ProductSrv) ProductDelete(ctx context.Context, req *types.ProductDeleteReq) (resp interface{}, err error) {
	u, _ := ctl.GetUserInfo(ctx)
	err = dao.NewProductDao(ctx).DeleteProduct(req.ID, u.Id)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}
	return
}

// 更新商品
func (s *ProductSrv) ProductUpdate(ctx context.Context, req *types.ProductUpdateReq) (resp interface{}, err error) {
	product := &model.Product{
		Name:       req.Name,
		CategoryID: req.CategoryID,
		Title:      req.Title,
		Info:       req.Info,
		// ImgPath:       service.ImgPath,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		OnSale:        req.OnSale,
	}
	err = dao.NewProductDao(ctx).UpdateProduct(req.ID, product)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}

	return
}

// 搜索商品 TODO 后续用脚本同步数据MySQL到ES，用ES进行搜索
func (s *ProductSrv) ProductSearch(ctx context.Context, req *types.ProductSearchReq) (resp interface{}, err error) {
	products, count, err := dao.NewProductDao(ctx).SearchProduct(req.Info, req.BasePage)
	if err != nil {
		log.LogrusObj.Error(err)
		return
	}

	pRespList := make([]*types.ProductResp, 0)
	for _, p := range products {
		pResp := &types.ProductResp{
			ID:            p.ID,
			Name:          p.Name,
			CategoryID:    p.CategoryID,
			Title:         p.Title,
			Info:          p.Info,
			ImgPath:       p.ImgPath,
			Price:         p.Price,
			DiscountPrice: p.DiscountPrice,
			View:          p.View(),
			CreatedAt:     p.CreatedAt.Unix(),
			Num:           p.Num,
			OnSale:        p.OnSale,
			BossID:        p.BossID,
			BossName:      p.BossName,
			BossAvatar:    p.BossAvatar,
		}
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			pResp.BossAvatar = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.AvatarPath + pResp.BossAvatar
			pResp.ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + pResp.ImgPath
		}
		pRespList = append(pRespList, pResp)
	}

	resp = &types.DataListResp{
		Item:  pRespList,
		Total: count,
	}

	return
}

// ProductImgList 获取商品列表图片
func (s *ProductSrv) ProductImgList(ctx context.Context, req *types.ListProductImgReq) (resp interface{}, err error) {
	productImgs, _ := dao.NewProductImgDao(ctx).ListProductImgByProductId(req.ID)
	for i := range productImgs {
		if conf.Config.System.UploadModel == consts.UploadModelLocal {
			productImgs[i].ImgPath = conf.Config.PhotoPath.PhotoHost + conf.Config.System.HttpPort + conf.Config.PhotoPath.ProductPath + productImgs[i].ImgPath
		}
	}

	resp = &types.DataListResp{
		Item:  productImgs,
		Total: int64(len(productImgs)),
	}

	return
}
