package v1

import (
	"errors"
	"g_mall/consts"
	"g_mall/pkg/utils/ctl"
	"g_mall/pkg/utils/log"
	"g_mall/service"
	"g_mall/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateProductHandler 创建商品
func CreateProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductCreateReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		form, _ := ctx.MultipartForm() //它通过 ctx.MultipartForm() 来获取前端随表单一起上传的图片文件。这是一个关键步骤，因为商品信息和图片是放在一个请求里（multipart/form-data类型）发过来的
		//       c.  调用服务层：将校验好的数据和图片文件，全部交给 service 层的 ProductCreate 方法去处理真正的业务逻辑。
		//          d.  返回响应：接收 service 层的处理结果，并将其封装成JSON格式返回给前端。
		files := form.File["image"]
		l := service.GetProductSrv()
		resp, err := l.ProductCreate(ctx.Request.Context(), files, &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// ListProductsHandler 商品列表
func ListProductsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductListReq
		//参数绑定
		if err := ctx.ShouldBind(&req); err != nil {
			// 2参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		if req.PageSize == 0 {
			req.PageSize = consts.BaseProductPageSize
		}

		l := service.GetProductSrv()
		//调用服务层
		resp, err := l.ProductList(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		//成功效应
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// ShowProductHandler 商品详情
func ShowProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductShowReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetProductSrv()
		resp, err := l.ProductShow(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// DeleteProductHandler 删除商品
func DeleteProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductDeleteReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetProductSrv()
		resp, err := l.ProductDelete(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// UpdateProductHandler 更新商品
func UpdateProductHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductUpdateReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetProductSrv()
		resp, err := l.ProductUpdate(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// SearchProductsHandler 搜索商品
func SearchProductsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ProductSearchReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		if req.PageSize == 0 {
			req.PageSize = consts.BasePageSize
		}

		l := service.GetProductSrv()
		resp, err := l.ProductSearch(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

func ListProductImgHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ListProductImgReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		if req.ID == 0 {
			err := errors.New("参数错误,id不能为空")
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetProductSrv()
		resp, err := l.ProductImgList(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}
