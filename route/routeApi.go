package route

import (
	"github.com/gin-gonic/gin"
	"mock/controller"
	"net/http"
	"strconv"
	"time"
)

/**
 * Created by zc on 2019-10-24.
 */
func routeApi(g *gin.Engine) {

	homeController := new(controller.HomeController)
	g.GET("/", homeController.Index)
	g.GET("/config", homeController.Config)

	userController := new(controller.UserController)
	g.GET("/user/list", userController.GetList)
	g.POST("/user/add", userController.CreateUser)
	g.POST("/tenant/add", userController.CreateTenant)
	g.POST("/tenant/reset", userController.ResetTenantIds)
	g.POST("/tenant/exchange", userController.ExchangeTenant)
	g.POST("/user/relate", userController.UserRelateTenant)
	g.POST("/user/mobileRest", userController.UserMobileRest)

	g.POST("/service/user/v1/get-token", userController.GetToken)             // (uuc)code换用户token
	g.POST("/service/api/v1/userinfo", userController.GetUserInfo)            // (uuc)用户token换用户信息
	g.GET("/service/api/v1/userinfo/tenement", userController.GetUserInfoAll) // (uuc)用户token换用户信息
	g.POST("/service/oauth/token", userController.Token)                      // (uuc)token处理
	g.POST("/cuba/getAllUser", userController.GetTenantAllUser)               // 获取租户下所有用户信息
	g.POST("/cuba/service/uuc_center/login", userController.LoginUserName)
	g.POST("/service/sso/v1/uuc/jump/uuc", userController.Uuc)
	
	g.GET("/01", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello 01!")
	})

	g.GET("/ts", func(c *gin.Context) {
		ts := time.Now().UnixNano() / 1e3
		c.String(http.StatusOK, strconv.Itoa(int(ts)))
	})

	mock := g.Group("/mock")
	{
		templateController := new(controller.TemplateController)
		mock.GET("/user/list", templateController.UserList)
		mock.GET("/api/list", templateController.ApiList)
		mock.GET("/contract", templateController.Contract)

		mockController := new(controller.MockController)
		mock.Any("/any", mockController.Any)
		mock.POST("/upload", mockController.Upload)

		storeController := new(controller.StoreController)
		mock.POST("/api/audit", storeController.AuditStatus)
		mock.POST("/contract", storeController.CreateContract)
	}
}
