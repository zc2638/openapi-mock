package data

import "strconv"

/**
 * Created by zc on 2019-10-24.
 */
type User struct {
	ID       string `json:"id"`       // 用户id
	UserName string `json:"username"` // 用户名称
	NickName string `json:"nickname"` // 用户昵称
	Phone    string `json:"phone"`    // 联系方式
	Gender   int    `json:"gender"`   // 性别
	HeadImg  string `json:"headImg"`  // 头像
	Code     string `json:"code"`     // 8位校验码
}

type UserData struct {
	User
	TenantList []UserTenantData `json:"tenementInfoList"` // 租户列表
}

type UserTenantData struct {
	ID       string `json:"id"`       // 租户id
	Name     string `json:"name"`     // 租户名称
	UserType int    `json:"userType"` // 用户类型：0自建 1授权
}

type UserToken struct {
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type AppToken struct {
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

type Tenant struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TenantSet []Tenant

func (s TenantSet) Len() int           { return len(s) }
func (s TenantSet) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TenantSet) Less(i, j int) bool {

	ii, _ := strconv.Atoi(s[i].ID)
	ij, _ := strconv.Atoi(s[j].ID)
	return ii < ij
}

type TenantData struct {
	Tenant
	UserList []TenantUserData `json:"userList"`
}

type TenantUserData struct {
	User
	UserType int `json:"userType"`
}

type TenantUserDataSet []TenantUserData

func (s TenantUserDataSet) Len() int           { return len(s) }
func (s TenantUserDataSet) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s TenantUserDataSet) Less(i, j int) bool {

	ii, _ := strconv.Atoi(s[i].ID)
	ij, _ := strconv.Atoi(s[j].ID)
	return ii < ij
}