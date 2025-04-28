package router

import (
	"github.com/gin-gonic/gin"
	"server/api"
	"server/middleware"
)

// UserRouter 用户路由
type UserRouter struct {
}

func (userRouter *UserRouter) InitUserRouter(privateRouter *gin.RouterGroup, publicRouter *gin.RouterGroup, adminRouter *gin.RouterGroup) {
	userPrivateRouter := privateRouter.Group("user")
	userPublicRouter := publicRouter.Group("user")
	userLoginRouter := publicRouter.Group("user").Use(middleware.LoginRecord())
	userAdminRouter := adminRouter.Group("user")

	userApi := api.ApiGroupApp.UserApi
	{
		userPrivateRouter.POST("logout", userApi.Logout)
		userPrivateRouter.PUT("resetPassword", userApi.UserResetPassword)
		userPrivateRouter.GET("info", userApi.UserInfo)
		userPrivateRouter.PUT("changeInfo", userApi.UserChangeInfo)
		userPrivateRouter.GET("weather", userApi.UserWeather)
		userPrivateRouter.GET("chart", userApi.UserChart)
	}
	{
		userPublicRouter.POST("forgotPassword", userApi.ForgotPassword)
		userPublicRouter.GET("card", userApi.UserCard)
	}
	{
		userLoginRouter.POST("register", userApi.Register)
		userLoginRouter.POST("login", userApi.Login)
	}
	{
		userAdminRouter.GET("list", userApi.UserList)
		//userAdminRouter.PUT("freeze", userApi.UserFreeze)
		//userAdminRouter.PUT("unfreeze", userApi.UserUnfreeze)
		//userAdminRouter.GET("loginList", userApi.UserLoginList)
	}
}
