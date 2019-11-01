package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/**
 * Created by zc on 2019-10-24.
 */
type BaseController struct {}

func (t *BaseController) Api(c *gin.Context, code int, data interface{}) {
	c.JSON(code, data)
}

func (t *BaseController) Succ(c *gin.Context, msg string) {
	t.Api(c, http.StatusOK, gin.H{
		"status": "success",
		"message": msg,
	})
}

func (t *BaseController) Data(c *gin.Context, data interface{}) {
	t.Api(c, http.StatusOK, data)
}

func (t *BaseController) Err(c *gin.Context, err error) {

	fmt.Println("[Error]", err.Error())
	t.Api(c, http.StatusBadRequest, gin.H{
		"msg": err.Error(),
	})
}

func (t *BaseController) ErrData(c *gin.Context, err error) {
	t.Api(c, http.StatusOK, gin.H{
		"status": "error",
		"message": err.Error(),
	})
}

const (
	ErrorRequest = RequestError("请求异常")
	AuthCodeError = RequestError("authCode错误")
	TokenError = RequestError("token异常")
	AuthError = RequestError("身份认证失败")
	TenantError = RequestError("租户认证失败")
	TenantRepeat = RequestError("租户已存在")
)

const (
	StatusApply    int = iota // 申请中
	StatusRefuse              // 拒绝
	StatusApproval            // 通过
)

type RequestError string

func (e RequestError) Error() string { return string(e) }