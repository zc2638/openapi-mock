package controller

import (
	"github.com/gin-gonic/gin"
	"mock/service"
	"mock/template"
	"net/http"
)

/**
 * Created by zc on 2019-11-01.
 */
type TemplateController struct{ BaseController }

func (t *TemplateController) UserList(c *gin.Context) {

	userService := service.UserService{}
	users, err := userService.GetUsers()
	if err != nil {
		c.String(http.StatusOK, template.Error(err))
		return
	}

	tenants, err := userService.GetTenants()
	if err != nil {
		c.String(http.StatusOK, template.Error(err))
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, template.UserListTemplate(tenants, users))
}

func (t *TemplateController) ApiList(c *gin.Context) {

	storeService := service.StoreService{}
	set, err := storeService.GetApiSet()
	if err != nil {
		c.String(http.StatusOK, template.Error(err))
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, template.ApiListTemplate(set))
}

func (t *TemplateController) Contract(c *gin.Context) {

	userService := service.UserService{}
	users, err := userService.GetUsers()
	if err != nil {
		c.String(http.StatusOK, template.Error(err))
		return
	}

	storeService := service.StoreService{}
	set, err := storeService.GetApiSet()
	if err != nil {
		c.String(http.StatusOK, template.Error(err))
		return
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, template.ContractTemplate(users, set))
}