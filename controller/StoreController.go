package controller

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/zc2638/go-validator"
	"mock/data"
	"mock/service"
	"strconv"
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
		"status": "success",
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
		"status": "success",
		"msg":    "下架成功",
	})
}

// (主动)修改审核状态
func (t *StoreController) AuditStatus(c *gin.Context) {

	apiId := c.PostForm("apiId")
	status := c.PostForm("status")

	validate := validator.NewVdr()
	validate.MakeValue(apiId, "required", "msg=缺少api id参数")
	validate.MakeValue(status, "reg=1|2|3", "msg=修改状态异常")
	if err := validate.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	sts, err := strconv.Atoi(status)
	if err != nil {
		t.ErrData(c, err)
		return
	}

	storeService := service.StoreService{}
	if err := storeService.Audit(apiId, sts); err != nil {
		t.ErrData(c, err)
		return
	}

	t.Data(c, gin.H{
		"status":  "success",
		"message": "操作成功",
	})
}

// (主动)创建服务合同
func (t *StoreController) CreateContract(c *gin.Context) {

	userId := c.PostForm("userId")
	tenantId := c.PostForm("tenantId")
	info := c.PostForm("info")

	validate := validator.NewVdr()
	validate.MakeValue(userId, "required", "msg=请选择购买用户")
	validate.MakeValue(tenantId, "required", "msg=请选择租户")
	if err := validate.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	var apiOfSCs []data.ApiOfSCs
	if err := json.Unmarshal([]byte(info), &apiOfSCs); err != nil {
		t.ErrData(c, errors.New("api解析失败"))
		return
	}

	if len(apiOfSCs) == 0 {
		t.ErrData(c, errors.New("请选择api"))
		return
	}

	storeService := service.StoreService{}
	if err := storeService.Contract(userId, tenantId, apiOfSCs); err != nil {
		t.ErrData(c, err)
		return
	}

	t.Data(c, gin.H{
		"status":  "success",
		"message": "操作成功",
	})
}
