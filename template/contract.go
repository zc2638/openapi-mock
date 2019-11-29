package template

import (
	"encoding/json"
	"mock/data"
)

/**
 * Created by zc on 2019-11-25.
 */
func ContractTemplate(users []data.UserData, apiSet []data.ApiData) string {

	userList := "<option value=''>请选择用户</option>"
	for _, u := range users {
		ub, _ := json.Marshal(u)
		userList += `<option value='` + string(ub) + `'>` + u.UserName + `</option>`
	}

	apiList := ""
	for _, api := range apiSet {
		if api.Status != 1 {
			continue
		}
		apiList += `<option value="` + api.ApiId + `">` + api.ApiName + "|" + api.TenantName + "|" + api.UserName + `</option>`
	}

	return `<!doctype html>
<html xmlns=http://www.w3.org/1999/xhtml>
<meta charset=utf-8>
<title>OpenAPI-mock Contract</title>
<style>
    .item {border: 1px solid #ccc; padding: 20px 10px; margin-bottom: 20px;}
</style>
<body>
<div>
	<input type="button" name="submit" onclick="showApi()" value="API列表"
		   style="width: 80px;height: 20px; background: #ccc; border-radius: 5px;color: blue;margin-bottom: 20px;cursor: pointer;" />
</div>
<div id="contract" class="contract">
    <label id="user" style="margin: 0 20px;display: block;">
        <span>用户：</span>
        <select name="userInfo" style="height: 30px; width: 300px;" onchange="show()">
            ` + userList + `
        </select>
    </label><br/>
    <label id="tenant" style="margin: 0 20px;display: block"></label><br/>
    <label style="margin: 0 20px;display: block">
        <input type="button" id="add" name="add" value="添加API" onclick="add()"
               style="width: 80px;height: 20px; background: #ccc; border-radius: 5px;cursor: pointer;" />
    </label><br/>
    <div id="set"></div>
    <label style="margin: 0 20px;display: block">
        <input type="submit" name="submit" onclick="contractSubmit()"
               style="width: 80px;height: 20px; background: #ccc; border-radius: 5px;cursor: pointer;" />
    </label>
</div>
</body>
<script type="text/template" id="itemTemplate">
    <label style="margin: 0 20px;display: block">
        <span style="margin-top: 20px">api：</span>
        <select name="api" style="height: 30px; width: 300px;">
            ` + apiList + `
        </select>
    </label><br/>
    <label style="margin: 0 20px;display: block">
        <span>次数：</span>
        <input name="limit" value="0">
    </label>
    <label style="margin: 0 20px;display: block">
        <span>每秒钟限流：</span>
        <input name="second" value="0">
    </label>
    <label style="margin: 0 20px;display: block">
        <span>每分钟限流：</span>
        <input name="minute" value="0">
    </label>
    <label style="margin: 0 20px;display: block">
        <span>每小时限流：</span>
        <input name="hour" value="0">
    </label>
    <label style="margin: 0 20px;display: block">
        <span>每天限流：</span>
        <input name="day" value="0">
    </label>
    <label style="margin: 10px 20px 0 20px;display: block">
        <input name="del" type="button" value="删除" onclick="remove(this)" style="background-color: red; border-radius: 8px; width: 80px; height: 30px;cursor: pointer;color: #fff;">
    </label>
</script>
<script>
	const host = "http://" + window.location.host
	function showApi() {
		window.location.href = host + "/mock/api/list";
	}
	
    function show() {
        const userInfoStr = document.getElementsByName("userInfo")[0].value;
        const userInfo = JSON.parse(userInfoStr);
        if (!userInfo.hasOwnProperty("tenementInfoList") || !userInfo.tenementInfoList) {
            alert("请注意：该用户未关联任何租户");
            return;
        }
        let str = "";
        for (let t of userInfo.tenementInfoList) {
            str += "<option value=\"" + t.id +  "\">" + t.name + "</option>";
        }
        document.getElementById("tenant").innerHTML = "<span style=\"margin-top: 20px\">租户：</span><select name=\"tenantId\" style=\"height: 30px; width: 300px;\">" + str + "</select>";
    }
    
    function add() {
        let div = document.createElement("div");
        div.className = "item";
        div.innerHTML = document.getElementById("itemTemplate").innerHTML;
        document.getElementById("set").appendChild(div);
    }

    function remove(e) {
        document.getElementById("set").removeChild(e.parentElement.parentElement)
    }

    function contractSubmit() {
        const apiSet = document.getElementsByName("api");
        const limitSet = document.getElementsByName("limit");
        const secondSet = document.getElementsByName("second");
        const minuteSet = document.getElementsByName("minute");
        const hourSet = document.getElementsByName("hour");
        const daySet = document.getElementsByName("day");
        const apiNum = apiSet.length;
        if (apiNum !== limitSet.length || apiNum !== secondSet.length || apiNum !== minuteSet.length || apiNum !== hourSet.length || apiNum !== daySet.length) {
            alert("参数异常");
            return;
        }

        let info = [];
        for (let i = 0; i < apiNum; i++) {
			const api = apiSet[i].value;
			if (api === "") {
				alert("有未选择的API")
				return;
			}
            const data = {
                idApi: api,
                totalTimes: parseInt(limitSet[i].value),
                tenantEverySecond: parseInt(secondSet[i].value),
                tenantEveryMinute: parseInt(minuteSet[i].value),
                tenantEveryHour: parseInt(hourSet[i].value),
                tenantEveryDay: parseInt(daySet[i].value),
            };
            info.push(data);
        }

        const tenantId = document.getElementsByName("tenantId")[0].value;
        const userInfoStr = document.getElementsByName("userInfo")[0].value;
		const userInfo = JSON.parse(userInfoStr);
        const userId = userInfo.id;
        const data = "tenantId=" + tenantId + "&userId=" + userId + "&info=" + JSON.stringify(info);
        createRequest(host + "/mock/contract", "POST", data, function (res) {
            let result = JSON.parse(res);
            alert(result.message || "请求异常");
            window.location.reload();
        })
    }

    function createRequest(host, method, data, callback) {
        let xhr = window.XMLHttpRequest ? new window.XMLHttpRequest() :
            new window.ActiveXObject('Microsoft.XMLHTTP');
        xhr.open(method || "GET", host, false);
        xhr.onreadystatechange = function () {
            //判断请求状态是否是已经完成
            if (xhr.readyState === 4) {
                //判断服务器是否返回成功200,304
                if (xhr.status >= 200 && xhr.status <= 300 || xhr.status === 304) {
                    //接收xhr的数据
                    callback(xhr.responseText);
                }
            }
        };
        xhr.setRequestHeader("Content-Type", "application/x-www-form-urlencoded");
        xhr.send(data);
    }
</script>
</html>`
}