package listener

import (
	"encoding/json"
	"github.com/zc2638/go-validator"
	"mock/data"
	"mock/service"
	"mock/lib/rabbit"
)

/**
 * Created by zc on 2019-11-28.
 */
// 申请上架store
func storeApply() {
	ch := rabbit.Channel{}
	set, err := ch.Receive("store.apply")
	if err != nil {
		errorf("fail to register a consume(store.apply): %s", err)
		return
	}

	for d := range set {
		printf("Received a message: %s", d.Body)

		var apiData data.ApiData
		if err := json.Unmarshal(d.Body, &apiData); err != nil {
			errorf("marshal error: %s", err)
			continue
		}

		validate := validator.NewVdr().MakeStruct(apiData)
		if err := validate.Check(); err != nil {
			errorln("valid error: %s", err)
			continue
		}

		storeService := service.StoreService{}
		if err := storeService.Apply(apiData); err != nil {
			errorf("exec error: %s", err)
			continue
		}

	}
}

// 强制从store下架
func storeForce() {
	ch := rabbit.Channel{}
	set, err := ch.Receive("store.force")
	if err != nil {
		errorf("fail to register a consume(store.force): %s", err)
		return
	}

	for d := range set {
		printf("Received a message: %s", d.Body)

		var info map[string]string
		if err := json.Unmarshal(d.Body, &info); err != nil {
			errorf("marshal error: %s", err)
			continue
		}

		apiId, ok := info["apiId"]
		if !ok || apiId == "" {
			errorln("valid error: 缺少apiId参数")
			continue
		}

		storeService := service.StoreService{}
		if err := storeService.Force(apiId); err != nil {
			errorf("exec error: %s", err)
			continue
		}
	}
}