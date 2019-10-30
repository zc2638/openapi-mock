package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zc2638/go-validator"
	"mock/data"
	"mock/util/db"
)

/**
 * Created by zc on 2019-10-24.
 */
type StoreController struct{ BaseController }

// 申请上架
func (t *StoreController) Apply(c *gin.Context) {

	tenantName := c.PostForm("tenantName")
	userName := c.PostForm("userName")
	apiId := c.PostForm("apiId")
	apiName := c.PostForm("apiName")
	apiDesc := c.PostForm("apiDesc")

	validate := validator.NewVdr().
		MakeValue(tenantName, "required", "msg=缺少租户名称").
		MakeValue(userName, "required", "msg=缺少用户名称").
		MakeValue(apiId, "required", "msg=缺少api id参数").
		MakeValue(apiName, "required", "msg=缺少api名称").
		MakeValue(apiDesc, "required", "msg=缺少api描述")
	if err := validate.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	b, err := json.Marshal(data.ApiData{
		TenantName: tenantName,
		UserName:   userName,
		ApiId:      apiId,
		ApiName:    apiName,
		ApiDesc:    apiDesc,
		Status:     StatusApply,
	})
	if err != nil {
		t.ErrData(c, err)
		return
	}

	if err := db.Update(Store, apiId, string(b)); err != nil {
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

	_, err := db.View(Store, apiId)
	if err != nil {
		t.ErrData(c, err)
		return
	}

	if err := db.Delete(Store, apiId); err != nil {
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
