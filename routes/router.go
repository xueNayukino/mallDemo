package routes

import (
	api "g_mall/api/v1"
	"g_mall/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter() *gin.Engine {
	r := gin.Default()
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(middleware.Cors(), middleware.Jaeger())
	r.Use(sessions.Sessions("mysession", store))
	r.StaticFS("/static", http.Dir("./static"))
	v1 := r.Group("api/v1")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "success")
		})
		// 用户操作
		v1.POST("user/register", api.UserRegisterHandler())
		v1.POST("user/login", api.UserLoginHandler())

		// 商品操作
		v1.GET("product/list", api.ListProductsHandler())
		v1.GET("product/show", api.ShowProductHandler())
		v1.POST("product/search", api.SearchProductsHandler())
		v1.GET("product/imgs/list", api.ListProductImgHandler()) // 商品图片
		v1.GET("category/list", api.ListCategoryHandler())       // 商品分类
		v1.GET("carousels", api.ListCarouselsHandler())          // 轮播图

		authed := v1.Group("/") // 需要登陆保护
		authed.Use(middleware.AuthMiddleware())
		{

			// 用户操作
			authed.POST("user/update", api.UserUpdateHandler())
			authed.GET("user/show_info", api.ShowUserInfoHandler())
			authed.POST("user/send_email", api.SendEmailHandler())
			authed.GET("user/valid_email", api.ValidEmailHandler())
			authed.POST("user/following", api.UserFollowingHandler())
			authed.POST("user/unfollowing", api.UserUnFollowingHandler())
			authed.POST("user/avatar", api.UploadAvatarHandler()) // 上传头像

			// 商品操作
			authed.POST("product/create", api.CreateProductHandler())
			authed.POST("product/update", api.UpdateProductHandler())
			authed.POST("product/delete", api.DeleteProductHandler())

			// 收藏夹，无更新操作
			authed.GET("favorites/list", api.ListFavoritesHandler())
			authed.POST("favorites/create", api.CreateFavoriteHandler())
			authed.POST("favorites/delete", api.DeleteFavoriteHandler())
			// 订单操作
			authed.POST("orders/create", api.CreateOrderHandler())
			authed.GET("orders/list", api.ListOrdersHandler())
			authed.GET("orders/show", api.ShowOrderHandler())
			authed.POST("orders/delete", api.DeleteOrderHandler())

			// 购物车
			authed.POST("carts/create", api.CreateCartHandler())
			authed.GET("carts/list", api.ListCartHandler())
			authed.POST("carts/update", api.UpdateCartHandler()) // 购物车id
			authed.POST("carts/delete", api.DeleteCartHandler())

			// 收获地址操作
			authed.POST("addresses/create", api.CreateAddressHandler())
			authed.GET("addresses/show", api.ShowAddressHandler())
			authed.GET("addresses/list", api.ListAddressHandler())
			authed.POST("addresses/update", api.UpdateAddressHandler())
			authed.POST("addresses/delete", api.DeleteAddressHandler())

			// 支付功能
			authed.POST("paydown", api.OrderPaymentHandler())

			// 显示金额
			authed.POST("money", api.ShowMoneyHandler())

			// 秒杀专场
			authed.POST("import_skill_goods", api.ListSkillProductHandler())
			authed.POST("init_skill_goods", api.InitSkillProductHandler())
			authed.POST("skill_goods", api.SkillProductHandler())
		}
	}
	return r
}
