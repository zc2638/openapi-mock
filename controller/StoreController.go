package controller

import "github.com/gin-gonic/gin"

/**
 * Created by zc on 2019-10-24.
 */
type StoreController struct{ BaseController }

// 申请上架
func (t *StoreController) Apply(c *gin.Context) {

}

// 强制下架
func (t *StoreController) Force(c *gin.Context) {

}

// (主动)修改审核状态
func (t *StoreController) AuditStatus(c *gin.Context) {

}

// (主动)创建服务合同
func (t *StoreController) CreateContract(c *gin.Context) {

}
