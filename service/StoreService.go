package service

import (
	"encoding/json"
	"mock/data"
	"mock/lib/db"
	"mock/lib/rabbit"
	"strconv"
)

/**
 * Created by zc on 2019-10-31.
 */
type StoreService struct{ BaseService }

func (s *StoreService) GetApiSet() ([]data.ApiData, error) {

	b, err := db.View(db.Store, "apiSet")
	if err != nil {
		return nil, err
	}

	var set []data.ApiData
	if b == nil {
		return set, err
	}

	if err := json.Unmarshal(b, &set); err != nil {
		return nil, err
	}
	return set, err
}

func (s *StoreService) Apply(apiData data.ApiData) error {

	set, err := s.GetApiSet()
	if err != nil {
		return err
	}

	// 重复判断
	for _, api := range set {
		if api.ApiId == apiData.ApiId {
			return APIExist
		}
	}
	set = append(set, apiData)
	b, err := json.Marshal(set)
	if err != nil {
		return err
	}
	return db.Update(db.Store, "apiSet", string(b))
}

func (s *StoreService) Force(apiId string) error {

	set, err := s.GetApiSet()
	if err != nil {
		return err
	}

	var exist bool
	var newSet []data.ApiData
	for _, api := range set {
		if api.ApiId == apiId {
			exist = true
			continue
		}
		newSet = append(newSet, api)
	}
	if !exist {
		return APINotExist
	}

	b, err := json.Marshal(newSet)
	if err != nil {
		return err
	}
	return db.Update(db.Store, "apiSet", string(b))
}

func (s *StoreService) Audit(apiId string, sts int) error {

	set, err := s.GetApiSet()
	if err != nil {
		return err
	}

	var exist bool
	var newSet []data.ApiData
	for _, api := range set {
		if api.ApiId == apiId {
			api.Status = sts
			exist = true
		}
		newSet = append(newSet, api)
	}

	if !exist {
		return APINotExist
	}

	b, err := json.Marshal(newSet)
	if err != nil {
		return err
	}

	var status string
	switch sts {
	case 1:
		status = "online"
	case 2:
		status = "rejected"
	case 3:
		status = "offlineByStore"
	}

	ch := rabbit.Channel{}
	if err := ch.Send("store.direct", "store.audit", map[string]string{
		"apiId":   apiId,
		"status":  status,
		"comment": "",
	}); err != nil {
		return err
	}

	if sts == 2 || sts == 3 {
		return s.Force(apiId)
	}
	return db.Update(db.Store, "apiSet", string(b))
}

func (s *StoreService) Contract(userId, tenantId string, apiOfSCsData []data.ApiOfSCs) error {

	// 获取用户信息
	userService := UserService{}
	users, err := userService.GetUsers()
	if err != nil {
		return err
	}

	var user data.UserData
	for _, u := range users {
		if u.ID == userId {
			user = u
			break
		}
	}
	if user.ID == "" {
		return UserNotExist
	}

	// 获取租户信息
	var tenant data.UserTenantData
	for _, t := range user.TenantList {
		if t.ID == tenantId {
			tenant = t
			break
		}
	}
	if tenant.ID == "" {
		return TenantNotExist
	}

	// 获取api信息
	set, err := s.GetApiSet()
	if err != nil {
		return err
	}
	// api根据租户分类
	apiMap := make(map[string][]data.ApiOfSCs)
	for _, api := range set {
		for _, apiOfSCs := range apiOfSCsData {
			if api.ApiId == apiOfSCs.IdApi {
				if _, ok := apiMap[api.TenantId]; !ok {
					apiMap[api.TenantId] = make([]data.ApiOfSCs, 0)
				}
				apiMap[api.TenantId] = append(apiMap[api.TenantId], apiOfSCs)
			}
		}
	}

	// 请求创建合同服务
	for tid, apis := range apiMap {

		ch := rabbit.Channel{}
		if err := ch.Send("store.direct", "store.contract", map[string]interface{}{
			"email": "",
			"username": user.UserName,
			"phone": user.Phone,
			"gender": strconv.Itoa(user.Gender),
			"uucUserId": user.ID,
			"cubaTenantId": tenant.ID,
			"cubaTenantName": tenant.Name,
			"cubaUserType": strconv.Itoa(tenant.UserType),
			"providerTenantId": tid,
			"apiOfSCs": apis,
		}); err != nil {
			return err
		}
	}
	return nil
}
