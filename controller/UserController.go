package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zc2638/go-validator"
	"mock/data"
	"mock/service"
	"mock/util/db"
	"mock/util/jwtUtil"
	"strings"
)

/**
 * Created by zc on 2019-10-24.
 */
type UserController struct{ BaseController }

// code换用户token
func (t *UserController) GetToken(c *gin.Context) {

	authCode := c.PostForm("authCode")
	authorization := c.GetHeader("Authorization")

	validate := validator.NewVdr().
		MakeValue(authCode, "required", "msg=authCode不存在").
		MakeValue(authorization, "required", "msg=身份认证失败")
	if err := validate.Check(); err != nil {
		t.Err(c, err)
		return
	}

	userService := new(service.UserService)
	if userService.CheckAppToken(authorization) == false {
		t.Err(c, AuthError)
		return
	}

	users, err := userService.GetUsers()
	if err != nil {
		t.Err(c, err)
		return
	}

	var user data.User
	for _, u := range users {
		if authCode == u.Code {
			user = u
			break
		}
	}

	if user.ID == "" {
		t.Err(c, AuthCodeError)
		return
	}

	userToken, err := userService.CreateUserToken(user)
	if err != nil {
		t.Err(c, TokenError)
		return
	}
	t.Data(c, userToken)
}

// 用户token换用户信息
func (t *UserController) GetUserInfo(c *gin.Context) {

	authorization := c.GetHeader("Authorization")
	authSlice := strings.Split(authorization, "Bearer ")
	if len(authSlice) != 2 {
		t.Err(c, AuthError)
		return
	}

	jwtResult, err := jwtUtil.ParseInfo(authSlice[1], "")
	if err != nil {
		t.Err(c, TokenError)
		return
	}

	userId, ok := jwtResult["info"].(map[string]interface{})["id"]
	if !ok {
		t.Err(c, TokenError)
		return
	}

	userService := new(service.UserService)
	users, err := userService.GetUsers()
	if err != nil {
		t.Err(c, err)
		return
	}

	var user data.User
	for _, u := range users {
		if u.ID == userId {
			user = u
			break
		}
	}

	if user.ID == "" {
		t.Err(c, TokenError)
		return
	}
	t.Data(c, user)
}

// token处理
func (t *UserController) Token(c *gin.Context) {

	grantType := c.PostForm("grant_type")
	refreshToken := c.PostForm("refresh_token")

	userService := new(service.UserService)
	if grantType == "client_credentials" { // 应用token生成
		appToken, err := userService.CreateAppToken()
		if err != nil {
			t.Err(c, TokenError)
			return
		}
		t.Data(c, appToken)
	} else if grantType == "refresh_token" { // 用户refreshToken换token
		userToken, err := userService.ParseRefreshToken(refreshToken)
		if err != nil {
			t.Err(c, TokenError)
			return
		}
		t.Data(c, userToken)
	} else {
		t.Err(c, ErrorRequest)
	}
}

// 校验用户是否在租户下
func (t *UserController) CheckTenantUser(c *gin.Context) {

	userId := c.PostForm("userId")
	tenantId := c.PostForm("tenantId")

	validate := validator.NewVdr().
		MakeValue(userId, "required", "msg=未找到用户").
		MakeValue(tenantId, "required", "msg=未找到租户")
	if err := validate.Check(); err != nil {
		t.Err(c, err)
		return
	}

	userService := new(service.UserService)
	tenantInfo := userService.CheckTenantUser(userId, tenantId)
	if tenantInfo.ID == "" {
		t.Err(c, AuthError)
		return
	}
	t.Data(c, tenantInfo)
}

// 获取租户下所有用户信息
func (t *UserController) GetTenantAllUser(c *gin.Context) {

	tenantId := c.PostForm("tenantId")
	validate := validator.NewVdr().MakeValue(tenantId, "required", "msg=未找到租户")
	if err := validate.Check(); err != nil {
		t.Err(c, err)
		return
	}

	userService := new(service.UserService)
	tenants, err := userService.GetTenants()
	if err != nil {
		t.Err(c, err)
		return
	}

	var tenant data.Tenant
	for _, t := range tenants {
		if t.ID == tenantId {
			tenant = t
			break
		}
	}
	if tenant.ID == "" {
		t.Err(c, TenantError)
		return
	}

	users, err := userService.GetUsers()
	if err != nil {
		t.Err(c, err)
	}

	tenantUserData := make([]data.TenantUserInfo, 0)
	for _, u := range users {
		for _, tu := range tenant.Users {
			if u.ID == tu.ID {
				tenantUserData = append(tenantUserData, data.TenantUserInfo{
					User: u,
					Role: tu.Role,
				})
			}
		}
	}
	t.Data(c, gin.H{
		"userList": tenantUserData,
	})
}

// 修改用户的openAPI角色
func (t *UserController) ChangeUserRole(c *gin.Context) {

	userId := c.PostForm("userId")
	tenantId := c.PostForm("tenantId")
	role := c.PostForm("role")

	validate := validator.NewVdr()
	validate.MakeValue(userId, "required", "msg=未找到用户")
	validate.MakeValue(tenantId, "required", "msg=未找到租户")
	validate.MakeValue(role, "required", "msg=角色不能为空")
	if err := validate.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	userService := new(service.UserService)
	tenantInfo := userService.CheckTenantUser(userId, tenantId)
	if tenantInfo.ID == "" {
		t.ErrData(c, AuthError)
		return
	}

	if err := db.Update(CUBA, tenantId + userId, role); err != nil {
		t.ErrData(c, err)
		return
	}
	t.Data(c, gin.H{
		"status": "success",
		"message": "操作成功",
	})
}