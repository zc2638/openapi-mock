package service

import (
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"mock/data"
	"mock/util/db"
	"mock/util/jwtUtil"
	"strconv"
	"strings"
	"time"
)

/**
 * Created by zc on 2019-10-24.
 */
type UserService struct{ BaseService }

// 获取用户
func (s *UserService) GetUsers() ([]data.UserData, error) {

	bu, err := db.View(db.UUC, "user")
	if err != nil {
		return nil, err
	}

	var users []data.UserData
	if bu == nil {
		return users, err
	}

	if err := json.Unmarshal(bu, &users); err != nil {
		return nil, err
	}
	return users, err
}

// 获取租户信息
func (s *UserService) GetTenants() ([]data.Tenant, error) {

	b, err := db.View(db.CUBA, "tenant")
	if err != nil {
		return nil, err
	}

	var tenants []data.Tenant
	if b == nil {
		return tenants, err
	}

	if err := json.Unmarshal(b, &tenants); err != nil {
		return nil, err
	}
	return tenants, err
}

// 添加租户
func (s *UserService) CreateTenant(name string) error {

	tenants, err := s.GetTenants()
	if err != nil {
		return err
	}

	for _, tenant := range tenants {
		if tenant.Name == name {
			return TenantRepeat
		}
	}

	tenants = append(tenants, data.Tenant{
		ID:   uuid.NewV4().String(),
		Name: name,
	})

	bt, err := json.Marshal(tenants)
	if err != nil {
		return err
	}
	return db.Update(db.CUBA, "tenant", string(bt))
}

// 添加用户
func (s *UserService) CreateUser(username, nickname, phone string) error {

	users, err := s.GetUsers()
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.UserName == username {
			return UserRepeat
		}
		if user.Phone == phone {
			return PhoneRepeat
		}
	}

	num := len(users) + 10000001
	users = append(users, data.UserData{
		User: data.User{
			ID:       uuid.NewV4().String(),
			UserName: username,
			NickName: nickname,
			Phone:    phone,
			HeadImg:  "http://idiy.haier.com/upload/test/eeb3c64b-706c-4bb8-a8c3-dd25c6d35824.jpg",
			Code:     strconv.Itoa(num),
			Gender:   2,
		},
	})

	bu, err := json.Marshal(users)
	if err != nil {
		return err
	}
	fmt.Println(string(bu))
	return db.Update(db.UUC, "user", string(bu))
}

// 用户关联租户
func (s *UserService) UserRelateTenant(userId, tenantId, userType string) error {

	tenants, err := s.GetTenants()
	if err != nil {
		return err
	}

	var tenant data.Tenant
	for _, t := range tenants {
		if t.ID == tenantId {
			tenant = t
			break
		}
	}
	if tenant.ID == "" {
		return TenantNotExist
	}

	users, err := s.GetUsers()
	if err != nil {
		return err
	}

	userTypeInt, err := strconv.Atoi(userType)
	if err != nil {
		return err
	}

	var newUsers []data.UserData
	for _, u := range users {
		user := u
		if u.ID == userId {
			tenantList := make([]data.UserTenantData, 0)
			if user.TenantList != nil {
				for _, td := range user.TenantList {
					if td.ID == tenantId {
						return UserRelateRepeat
					}
				}
				tenantList = user.TenantList
			}
			tenantList = append(tenantList, data.UserTenantData{
				ID:       tenant.ID,
				Name:     tenant.Name,
				UserType: userTypeInt,
			})
			user.TenantList = tenantList
		}
		newUsers = append(newUsers, user)
	}

	bu, err := json.Marshal(newUsers)
	if err != nil {
		return err
	}
	return db.Update(db.UUC, "user", string(bu))
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
		Scope:        "read write",
		TokenType:    "bearer",
		ExpiresIn:    int64(expireIn / 1000 / 1000),
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

// 生成应用token
func (s *UserService) CreateAppToken() (data.AppToken, error) {

	expireIn := time.Hour * 24
	appToken, err := jwtUtil.Create(nil, "openAPI", time.Now().Add(expireIn).Unix())

	return data.AppToken{
		Scope:       "read write",
		TokenType:   "bearer",
		ExpiresIn:   int64(expireIn / 1000 / 1000),
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
			user = u.User
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
