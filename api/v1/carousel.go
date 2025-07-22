package v1

import (
	"g_mall/pkg/utils/ctl"
	"g_mall/pkg/utils/log"
	"g_mall/service"
	"g_mall/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ListCarouselsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req types.ListCarouselReq
		if err := ctx.ShouldBind(&req); err != nil {
			// 参数校验
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}

		l := service.GetCarouselSrv()
		resp, err := l.ListCarousel(ctx.Request.Context(), &req)
		if err != nil {
			log.LogrusObj.Infoln(err)
			ctx.JSON(http.StatusOK, ErrorResponse(ctx, err))
			return
		}
		ctx.JSON(http.StatusOK, ctl.RespSuccess(ctx, resp))
	}
}
