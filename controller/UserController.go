package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/zc2638/go-validator"
	"mock/data"
	"mock/lib/jwtUtil"
	"mock/service"
	"net/url"
	"strings"
)

/**
 * Created by zc on 2019-10-24.
 */
type UserController struct{ BaseController }

// 获取所有租户和用户信息
func (t *UserController) GetList(c *gin.Context) {

	userService := service.UserService{}
	tenants, err := userService.GetTenants()
	if err != nil {
		t.Err(c, err)
		return
	}

	users, err := userService.GetUsers()
	if err != nil {
		t.Err(c, err)
		return
	}
	t.Data(c, gin.H{
		"tenants": tenants,
		"users":   users,
	})
}

// 新增租户
func (t *UserController) CreateTenant(c *gin.Context) {

	name := c.PostForm("name")
	if err := validator.NewVdr().MakeValue(name, "required", "msg=租户名称不能为空").Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	userService := service.UserService{}
	if err := userService.CreateTenant(name); err != nil {
		t.ErrData(c, err)
		return
	}
	t.Succ(c, "操作成功")
}

// 重置租户id
func (t *UserController) ResetTenantIds(c *gin.Context) {

	userService := service.UserService{}
	if err := userService.ChangeTenantIds(); err != nil {
		t.ErrData(c, err)
		return
	}
	t.Succ(c, "操作成功")
}

// 随意更换租户id
func (t *UserController) ExchangeTenant(c *gin.Context) {

	id := c.PostForm("id")
	exchangeId := c.PostForm("exchangeId")
	vdr := validator.NewVdr()
	vdr.MakeValue(id, `reg=^[1-9]\d*$`, "msg=id错误")
	vdr.MakeValue(exchangeId, `reg=^[1-9]\d*$`, "msg=exchangeId错误")
	if err := vdr.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	userService := service.UserService{}
	if err := userService.ExchangeTenant(id, exchangeId); err != nil {
		t.ErrData(c, err)
		return
	}
	t.Succ(c, "操作成功")
}

// 新增用户
func (t *UserController) CreateUser(c *gin.Context) {

	username := c.PostForm("username")
	nickname := c.PostForm("nickname")
	phone := c.PostForm("phone")

	vdr := validator.NewVdr()
	vdr.MakeValue(username, "required", "msg=请填写用户名")
	vdr.MakeValue(nickname, "required", "msg=请填写昵称")
	vdr.MakeValue(phone, "reg=^[1]([3-9])[0-9]{9}$", "msg=请填写正确的联系方式")
	if err := vdr.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	userService := service.UserService{}
	if err := userService.CreateUser(username, nickname, phone); err != nil {
		t.ErrData(c, err)
		return
	}
	t.Succ(c, "操作成功")
}

// 用户关联租户
func (t *UserController) UserRelateTenant(c *gin.Context) {

	userId := c.PostForm("userId")
	tenantId := c.PostForm("tenantId")
	userType := c.PostForm("userType")

	vdr := validator.NewVdr()
	vdr.MakeValue(userId, "required", "msg=请选择用户")
	vdr.MakeValue(tenantId, "required", "msg=请选择租户")
	vdr.MakeValue(userType, "reg=0|1", "msg=请填写用户角色，0创建者，1用户")
	if err := vdr.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	userService := service.UserService{}
	if err := userService.UserRelateTenant(userId, tenantId, userType); err != nil {
		t.ErrData(c, err)
		return
	}
	t.Succ(c, "操作成功")
}

// 用户手机号重置
func (t *UserController) UserMobileRest(c *gin.Context) {
	userService := service.UserService{}
	if err := userService.RestUserMobile(); err != nil {
		t.ErrData(c, err)
		return
	}
	t.Succ(c, "操作成功")
}

// code换用户token
func (t *UserController) GetToken(c *gin.Context) {

	authCode := c.PostForm("authCode")
	authorization := c.GetHeader("Authorization")

	validate := validator.NewVdr().
		MakeValue(authCode, "required", "msg=authCode不存在").
		MakeValue(authorization, "required", "msg=身份认证失败")
	if err := validate.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	userService := new(service.UserService)
	//if userService.CheckAppToken(authorization) == false {
	//	t.Err(c, AuthError)
	//	return
	//}

	users, err := userService.GetUsers()
	if err != nil {
		t.ErrData(c, err)
		return
	}

	var user data.User
	for _, u := range users {
		if authCode == u.Code {
			user = u.User
			break
		}
	}

	if user.ID == "" {
		t.ErrData(c, AuthCodeError)
		return
	}

	userToken, err := userService.CreateUserToken(user)
	if err != nil {
		t.ErrData(c, TokenError)
		return
	}
	t.Data(c, userToken)
}

// 用户token换用户信息
func (t *UserController) GetUserInfo(c *gin.Context) {

	authorization := c.GetHeader("Authorization")
	authSlice := strings.Split(authorization, "Bearer ")
	if len(authSlice) != 2 {
		t.ErrData(c, AuthError)
		return
	}

	jwtResult, err := jwtUtil.ParseInfo(authSlice[1], "")
	if err != nil {
		t.ErrData(c, TokenError)
		return
	}

	userId, ok := jwtResult["info"].(map[string]interface{})["id"]
	if !ok {
		t.ErrData(c, TokenError)
		return
	}

	userService := new(service.UserService)
	users, err := userService.GetUsers()
	if err != nil {
		t.ErrData(c, err)
		return
	}

	var user data.User
	for _, u := range users {
		if u.ID == userId {
			user = u.User
			break
		}
	}

	if user.ID == "" {
		t.ErrData(c, TokenError)
		return
	}
	t.Data(c, user)
}

// 用户token获取用户信息及用户所有租户信息
func (t *UserController) GetUserInfoAll(c *gin.Context) {

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

	userService := service.UserService{}
	users, err := userService.GetUsers()
	if err != nil {
		t.ErrData(c, err)
		return
	}

	var userData data.UserData
	for _, user := range users {
		if user.ID == userId {
			userData = user
			break
		}
	}
	if userData.ID == "" {
		t.ErrData(c, service.UserNotExist)
		return
	}
	t.Data(c, userData)
}

// token处理
func (t *UserController) Token(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	userService := new(service.UserService)
	switch grantType {
	case "client_credentials": // 应用token生成
		appToken, err := userService.CreateAppToken()
		if err != nil {
			t.Err(c, TokenError)
			return
		}
		t.Data(c, appToken)
	case "refresh_token": // 用户refreshToken换token
		refreshToken := c.PostForm("refresh_token")
		userToken, err := userService.ParseRefreshToken(refreshToken)
		if err != nil {
			t.Err(c, TokenError)
			return
		}
		t.Data(c, userToken)
	case "password": // 用户名密码换token
		username := c.PostForm("username")
		// 不校验密码
		users, err := userService.GetUsers()
		if err != nil {
			t.ErrData(c, err)
			return
		}

		var user data.User
		for _, u := range users {
			if u.Phone == username {
				user = u.User
				break
			}
		}
		if user.ID == "" {
			t.ErrData(c, AuthCodeError)
			return
		}

		userToken, err := userService.CreateUserToken(user)
		if err != nil {
			t.ErrData(c, TokenError)
			return
		}
		t.Data(c, userToken)
	default:
		t.Err(c, ErrorRequest)
	}
}

// 获取租户下所有用户信息
func (t *UserController) GetTenantAllUser(c *gin.Context) {

	result := gin.H{
		"code": "0",
		"msg": "ok",
		"isSuccess": false,
		"imgHost": nil,
	}

	tenantId := c.PostForm("tenementId")
	validate := validator.NewVdr().MakeValue(tenantId, "required", "msg=未找到租户")
	if err := validate.Check(); err != nil {
		result["msg"] = err.Error()
		t.Data(c, result)
		return
	}

	userService := new(service.UserService)
	tenants, err := userService.GetTenants()
	if err != nil {
		result["msg"] = err.Error()
		t.Data(c, result)
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
		result["msg"] = TenantError.Error()
		t.Data(c, result)
		return
	}

	users, err := userService.GetUsers()
	if err != nil {
		result["msg"] = TenantError.Error()
		t.Data(c, result)
		return
	}

	tenantUserData := make([]data.TenantUserData, 0)
	for _, u := range users {
		if u.TenantList == nil {
			continue
		}
		for _, tl := range u.TenantList {
			if tl.ID == tenantId {
				tenantUserData = append(tenantUserData, data.TenantUserData{
					UserId:   u.ID,
					UserName: u.UserName,
					Phone:    u.Phone,
					HeadImg:  u.HeadImg,
					UserType: tl.UserType,
				})
			}
		}
	}

	result["isSuccess"] = true
	result["data"] = tenantUserData
	t.Data(c, result)
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

	//if err := db.Update(CUBA, tenantId+userId, role); err != nil {
	//	t.ErrData(c, err)
	//	return
	//}
	t.Data(c, gin.H{
		"status":  "success",
		"message": "操作成功",
	})
}

// CUBA用户名密码登陆接口
func (t *UserController) LoginUserName(c *gin.Context) {

	result := gin.H{
		"code": 400,
		"msg": "",
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	vdr := validator.NewVdr()
	vdr.MakeValue(username, "required", "msg=请输入用户名")
	vdr.MakeValue(password, "required", "msg=请输入密码")
	if err := vdr.Check(); err != nil {
		result["msg"] = err.Error()
		t.Data(c, result)
		return
	}

	userService := service.UserService{}
	user, err := userService.GetUserByUserName(username, password)
	if err != nil {
		result["msg"] = err.Error()
		t.Data(c, result)
		return
	}

	userToken, err := userService.CreateUserToken(user.User)
	if err != nil {
		result["msg"] = TokenError.Error()
		t.Data(c, result)
		return
	}

	t.Data(c, gin.H{
		"data": gin.H{
			"ucUserId": user.ID,
			"username": user.UserName,
			"ucUserToken": userToken.AccessToken,
			"nickname": user.NickName,
			"phone": user.Phone,
		},
		"code": 0,
		"msg": "",
	})
}

// UUC单点登录
func (t *UserController) Uuc(c *gin.Context) {
	authorization := c.GetHeader("Authorization")
	target := c.PostForm("target")
	targetUrl := c.PostForm("targetUrl")
	gateWayUrl := c.PostForm("gateWayUrl")
	state := c.PostForm("state")

	vdr := validator.NewVdr()
	vdr.MakeValue(authorization, "required", "msg=身份认证失败")
	vdr.MakeValue(target, "required", "msg=应用认证失败")
	vdr.MakeValue(targetUrl, "required", "msg=跳转地址不存在")
	vdr.MakeValue(gateWayUrl, "required", "msg=站点gateWay地址不存在")
	if err := vdr.Check(); err != nil {
		t.ErrData(c, err)
		return
	}

	userService := service.UserService{}
	userData, err := userService.ParseToken(authorization)
	if err != nil {
		t.ErrData(c, err)
		return
	}

	t.Data(c, gin.H{
		"authCode": userData.Code, // 授权码
		"target": "OpenAPI",
		"targetUrl": targetUrl,
		"redirectUrl": gateWayUrl + "?code=" + userData.Code + "&targetUrl=" + url.PathEscape(targetUrl) + "&state=" + state,
	})
}