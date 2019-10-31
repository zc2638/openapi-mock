package service

import (
	"encoding/json"
	"mock/data"
	"mock/util/db"
)

/**
 * Created by zc on 2019-10-31.
 */
type StoreService struct{ BaseService }

func (s *StoreService) Apply(apiData data.ApiData) error {

	b, err := json.Marshal(apiData)
	if err != nil {
		return err
	}
	return db.Update(db.Store, apiData.ApiId, string(b))
}

func (s *StoreService) Force(apiId string) error {
	return db.Delete(db.Store, apiId)
}
