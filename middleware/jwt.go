package middleware

import (
	"g_mall/consts"
	"g_mall/pkg/e"
	"g_mall/pkg/utils/ctl"
	util "g_mall/pkg/utils/jwt"
	"github.com/gin-gonic/gin"
	"strings"
)

// AuthMiddleware token验证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		code = e.SUCCESS
		accessToken := c.GetHeader("access_token")
		if strings.HasPrefix(accessToken, "Bearer ") {
			accessToken = accessToken[7:]
		}
		refreshToken := c.GetHeader("refresh_token")
		if accessToken == "" {
			code = e.InvalidParams
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "Token不能为空",
			})
			c.Abort()
			return
		}
		newAccessToken, newRefreshToken, err := util.ParseRefreshToken(accessToken, refreshToken)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
		}
		if code != e.SUCCESS {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   "鉴权失败",
				"error":  err.Error(),
			})
			c.Abort()
			return
		}
		claims, err := util.ParseToken(newAccessToken)
		if err != nil {
			code = e.ErrorAuthCheckTokenFail
			c.JSON(200, gin.H{
				"status": code,
				"msg":    e.GetMsg(code),
				"data":   err.Error(),
			})
			c.Abort()
			return
		}
		SetToken(c, newAccessToken, newRefreshToken)
		//下一行很重要
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.ID}))
		ctl.InitUserInfo(c.Request.Context())
		c.Next()
	}
}

func SetToken(c *gin.Context, accessToken, refreshToken string) {
	secure := IsHttps(c)
	c.Header(consts.AccessTokenHeader, accessToken)
	c.Header(consts.RefreshTokenHeader, refreshToken)
	c.SetCookie(consts.AccessTokenHeader, accessToken, consts.MaxAge, "/", "", secure, true)
	c.SetCookie(consts.RefreshTokenHeader, refreshToken, consts.MaxAge, "/", "", secure, true)
}

// 判断是否https
func IsHttps(c *gin.Context) bool {
	if c.GetHeader(consts.HeaderForwardedProto) == "https" || c.Request.TLS != nil {
		return true
	}
	return false
}
