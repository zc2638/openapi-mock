package controller

import (
	"github.com/gin-gonic/gin"
	"mock/config"
	"net/http"
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