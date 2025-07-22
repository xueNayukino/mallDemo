package v1

import (
	"g_mall/consts"
	"g_mall/pkg/utils/ctl"
	"g_mall/pkg/utils/log"
	"g_mall/service"
	"g_mall/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateCartHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.CartCreateReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetCartSrv()
		resp, err := l.CartCreate(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// ListCartHandler 购物车详细信息
func ListCartHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.CartListReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		if req.PageSize == 0 {
			req.PageSize = consts.BasePageSize
		}

		l := service.GetCartSrv()
		resp, err := l.CartList(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// UpdateCartHandler 修改购物车信息
func UpdateCartHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.UpdateCartServiceReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetCartSrv()
		resp, err := l.CartUpdate(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}

// DeleteCartHandler 删除购物车
func DeleteCartHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.CartDeleteReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetCartSrv()
		resp, err := l.CartDelete(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}
