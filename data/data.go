package data

/**
 * Created by zc on 2019-10-24.
 */
type User struct {
	ID       string `yaml:"id" json:"id"`
	UserName string `yaml:"userName" json:"username"`
	NickName string `yaml:"nickName" json:"nickname"`
	Phone    string `yaml:"phone" json:"phone"`
	Gender   int    `yaml:"gender"json:"gender"`
	HeadImg  string `yaml:"headImg" json:"headImg"`
	Code     string `yaml:"code" json:"code"`
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
	ID    string       `yaml:"id"`
	Name  string       `yaml:"name"`
	Desc  string       `yaml:"desc"`
	Users []TenantUser `yaml:"users"`
}

type TenantUser struct {
	ID   string `yaml:"id"`
	Role string `yaml:"role"`
}

type TenantInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
	Role string `json:"role"`
}

type TenantUserInfo struct {
	User
	Role string `json:"role"`
}
