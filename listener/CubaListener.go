package listener

import (
	"encoding/json"
	"mock/lib/rabbit"
)

/**
 * Created by zc on 2019-11-28.
 */
// 修改用户角色
func cubaRole() {
	ch := rabbit.Channel{}
	set, err := ch.Receive("cuba.role")
	if err != nil {
		errorf("fail to register a consume: %s", err)
		return
	}

	for d := range set {
		printf("Received a message: %s", d.Body)

		var info map[string]string
		if err := json.Unmarshal(d.Body, &info); err != nil {
			errorf("marshal error: %s", err)
			continue
		}

		tenantId, ok := info["tenantId"]
		if !ok || tenantId == "" {
			errorln("valid error: 缺少tenantId参数")
			continue
		}

		userId, ok := info["userId"]
		if !ok || userId == "" {
			errorln("valid error: 缺少userId参数")
			continue
		}

		role, ok := info["role"]
		if !ok || role == "" {
			errorln("valid error: 缺少role参数")
			continue
		}
	}
}