package controller

import (
	"github.com/gin-gonic/gin"
	"mock/config"
	"net/http"
	"strconv"
)

/**
 * Created by zc on 2019-10-24.
 */
type HomeController struct{ BaseController }

func (t *HomeController) Index(c *gin.Context) {
	c.String(http.StatusOK, "Hello World!")
}

func (t *HomeController) Config(c *gin.Context) {
	c.JSON(http.StatusOK, config.Cfg)
}

func (t *HomeController) HttpStatus(c *gin.Context) {
	statusCode := c.DefaultQuery("statusCode", "200")
	code, err := strconv.Atoi(statusCode)
	if err != nil {
		c.String(http.StatusOK, "解析失败：statusCode参数只支持http状态码整数")
	}
	c.String(code, "")
}