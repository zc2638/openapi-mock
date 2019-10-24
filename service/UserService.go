package service

import (
	"errors"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"mock/data"
	"mock/util/jwtUtil"
	"strings"
	"time"
)

/**
 * Created by zc on 2019-10-24.
 */
type UserService struct{ BaseService }

// 获取用户
func (s *UserService) GetUsers() ([]data.User, error) {

	b, err := ioutil.ReadFile("user.yml")
	if err != nil {
		return nil, err
	}

	var users []data.User
	err = yaml.Unmarshal(b, &users)
	return users, err
}

// 获取租户信息
func (s *UserService) GetTenants() ([]data.Tenant, error) {

	b, err := ioutil.ReadFile("tenant.yml")
	if err != nil {
		return nil, err
	}

	var tenants []data.Tenant
	err = yaml.Unmarshal(b, &tenants)
	return tenants, err
}

// 生成用户token
func (s *UserService) CreateUserToken(user data.User) (data.UserToken, error) {

	if user.ID == "" {
		return data.UserToken{}, errors.New("用户不存在")
	}

	// 生成用户token
	expireIn := time.Hour * 2
	token, err := jwtUtil.Create(map[string]interface{}{
		"id": user.ID,
	}, "", time.Now().Add(expireIn).Unix())
	if err != nil {
		return data.UserToken{}, err
	}

	// 生成用户refreshToken
	refreshToken, err := jwtUtil.Create(map[string]interface{}{
		"id": user.ID,
	}, "", time.Now().Add(time.Hour * 24).Unix())
	if err != nil {
		return data.UserToken{}, err
	}

	return data.UserToken{
		Scope: "read write",
		TokenType: "bearer",
		ExpiresIn: int64(expireIn / 1000 / 1000),
		AccessToken: token,
		RefreshToken: refreshToken,
	}, nil
}

// 生成应用token
func (s *UserService) CreateAppToken() (data.AppToken, error) {

	expireIn := time.Hour * 24
	appToken, err := jwtUtil.Create(nil, "openAPI", time.Now().Add(expireIn).Unix())

	return data.AppToken{
		Scope: "read write",
		TokenType: "bearer",
		ExpiresIn: int64(expireIn / 1000 / 1000),
		AccessToken: appToken,
	}, err
}

// 解析用户refreshToken
func (s *UserService) ParseRefreshToken(refreshToken string) (data.UserToken, error) {

	userToken := data.UserToken{}
	jwtResult, err := jwtUtil.ParseInfo(refreshToken, "")
	if err != nil {
		return userToken, err
	}

	userId, ok := jwtResult["info"].(map[string]interface{})["id"]
	if !ok {
		return userToken, errors.New("token异常")
	}

	users, err := s.GetUsers()
	if err != nil {
		return userToken, err
	}

	var user data.User
	for _, u := range users {
		if u.ID == userId {
			user = u
			break
		}
	}
	return s.CreateUserToken(user)
}

// 检查应用token
func (s *UserService) CheckAppToken(authorization string) bool {

	authSlice := strings.Split(authorization, "Bearer ")
	if len(authSlice) != 2 {
		return false
	}
	return jwtUtil.CheckValid(authSlice[1], "openAPI")
}
