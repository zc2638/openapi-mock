package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zc2638/go-validator"
	"mock/data"
	"mock/service"
)

/**
 * Created by zc on 2019-10-24.
 */
type StoreController struct{ BaseController }

// 申请上架
func (t *StoreController) Apply(c *gin.Context) {

	tenantId := c.PostForm("tenantId")
	tenantName := c.PostForm("tenantName")
	userId := c.PostForm("userId")
	userName := c.PostForm("userName")
	apiId := c.PostForm("apiId")
	apiName := c.PostForm("apiName")
	apiDesc := c.PostForm("apiDesc")

	validate := validator.NewVdr().
		MakeValue(tenantId, "required", "msg=缺少租户标识").
		MakeValue(tenantName, "required", "msg=缺少租户名称").
		MakeValue(userId, "required", "msg=缺少用户标识").
		MakeValue(userName, "required", "msg=缺少用户名称").
		MakeValue(apiId, "required", "msg=缺少api id参数").
		MakeValue(apiName, "required", "msg=缺少api名称").
		MakeValue(apiDesc, "required", "msg=缺少api描述")
	if err := validate.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	storeService := service.StoreService{}
	if err := storeService.Apply(data.ApiData{
		TenantId:   tenantId,
		TenantName: tenantName,
		UserId:     userId,
		UserName:   userName,
		ApiId:      apiId,
		ApiName:    apiName,
		ApiDesc:    apiDesc,
		Status:     StatusApply,
	}); err != nil {
		t.ErrData(c, err)
		return
	}

	t.Data(c, gin.H{
		"status": "ok",
		"msg":    "申请成功",
	})
}

// 强制下架
func (t *StoreController) Force(c *gin.Context) {

	apiId := c.PostForm("apiId")
	validate := validator.NewVdr().MakeValue(apiId, "required", "msg=缺少api id参数")
	if err := validate.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	storeService := service.StoreService{}
	if err := storeService.Force(apiId); err != nil {
		t.ErrData(c, err)
		return
	}

	t.Data(c, gin.H{
		"status": "ok",
		"msg":    "下架成功",
	})
}

// (主动)修改审核状态
func (t *StoreController) AuditStatus(c *gin.Context) {

}

// (主动)创建服务合同
func (t *StoreController) CreateContract(c *gin.Context) {

}
