package route

import (
	"github.com/gin-gonic/gin"
	"mock/controller"
)

/**
 * Created by zc on 2019-10-24.
 */
func routeApi(g *gin.Engine) {

	g.GET("/", new(controller.HomeController).Index)

	userController := new(controller.UserController)
	g.POST("/service/user/v1/get-token", userController.GetToken)  // (uuc)code换用户token
	g.POST("/service/api/v1/userinfo", userController.GetUserInfo) // (uuc)用户token换用户信息
	g.POST("/service/oauth/token", userController.Token)           // (uuc)token处理
	g.POST("/cuba/checkUser", userController.CheckTenantUser)      // 检查租户用户
	g.POST("/cuba/getAllUser", userController.GetTenantAllUser)    // 获取租户下所有用户信息

	storeController := new(controller.StoreController)
	g.POST("/store/apply", storeController.Apply)                   // (store)申请上架
	g.POST("/store/force", storeController.Force)                   // (store)强制下架
	g.POST("/store/auditStatus", storeController.AuditStatus)       // (store)修改审核状态（主动）
	g.POST("/store/createContract", storeController.CreateContract) // (store)创建服务合同（主动）
}
