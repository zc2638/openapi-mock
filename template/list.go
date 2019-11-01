package template

import (
	"mock/data"
	"strings"
)

/**
 * Created by zc on 2019-11-01.
 */
func UserListTemplate(tenants []data.Tenant, users []data.UserData) string {

	tenantList := ""
	if tenants != nil {
		for _, t := range tenants {
			tenantList += `<tr align="center">
                <td>` + t.ID + `</td>
                <td>` + t.Name + `</td>
            </tr>`
		}
	}

	userList := ""
	if users != nil {
		for _, u := range users {

			gender := "保密"
			switch u.Gender {
			case 0:
				gender = "女"
			case 1:
				gender = "男"
			}

			tenantName := ""
			if u.TenantList != nil {
				for _, t := range u.TenantList {
					tenantName += t.Name + `<br/>`
				}
				tenantName = strings.TrimRight(tenantName, "<br/>")
			}

			userList += `<tr align="center">
                <td>` + u.ID + `</td>
                <td>` + u.UserName + `</td>
                <td>` + u.NickName + `</td>
                <td>` + u.Phone + `</td>
                <td>` + gender + `</td>
                <td>` + u.Code + `</td>
                <td>` + tenantName + `</td>
                <td>
                    <button onclick="relate('` + u.ID + `')">关联租户</button>
                </td>
            </tr>`
		}
	}

	return `<!doctype html>
<html xmlns=http://www.w3.org/1999/xhtml>
<meta charset=utf-8>
<title>OpenAPI-mock用户中心</title>
<body>
<div class="content">
    <div class="pop"
         style="width: 100%; height: 100%; position: fixed; top: 0; left: 0; background: rgba(0,0,0,0.3);display: none;z-index: 100;"></div>
    <div class="tenantBox" style="position: absolute; left:50%; top:50%; display: none;z-index: 999;">
        <div class="item"
             style="width:550px; height:100px; position:absolute;left:-250px; top:-250px; border:2px solid; background: #fff;padding-top: 40px;">
            <label style="margin: 0 20px;display: block;">
                <span>租户名称：</span>
                <input name="tenantName" value="" placeholder="请输入租户名称" style="height: 30px; width: 300px;"/>
            </label><br/>
            <label style="margin: 0 20px;display: block">
                <input type="submit" name="submit" onclick="tenantSubmit()"
                       style="width: 80px;height: 20px; background: #ccc; border-radius: 5px;" />
            </label>
        </div>
    </div>
    <div class="userBox" style="position: absolute; left:50%; top:50%; display: none;z-index: 999;">
        <div class="item"
             style="width:550px; height:260px; position:absolute;left:-250px; top:-250px; border:2px solid; background: #fff;padding-top: 40px;">
            <label style="margin: 0 20px;display: block;">
                <span>用户名称：</span>
                <input name="userName" value="" placeholder="请输入用户名称" style="height: 30px; width: 300px;"/>
            </label><br/>
            <label style="margin: 0 20px;display: block">
                <span style="margin-top: 20px">用户昵称：</span>
                <input name="nickName" value="" placeholder="请输入用户昵称" style="height: 30px; width: 300px;"/>
            </label><br/>
            <label style="margin: 0 20px;display: block">
                <span style="margin-top: 20px">手机号码：</span>
                <input name="phone" value="" placeholder="请输入手机号码" style="height: 30px; width: 300px;"/>
            </label><br/>
            <label style="margin: 0 20px;display: block">
                <input type="submit" name="submit" onclick="userSubmit()"
                       style="width: 80px;height: 20px; background: #ccc; border-radius: 5px;" />
            </label>
        </div>
    </div>
    <div class="relateBox" style="position: absolute; left:50%; top:50%; display: none;z-index: 999;">
        <div class="item"
             style="width:550px; height:200px; position:absolute;left:-250px; top:-250px; border:2px solid; background: #fff;padding-top: 40px;">
            <label style="margin: 0 20px;display: block;">
                <span>租户id：</span>
                <input name="tenantId" value="" placeholder="请输入需要关联的租户id" style="height: 30px; width: 300px;"/>
            </label><br/>
            <label style="margin: 0 20px;display: block">
                <span style="margin-top: 20px">用户角色：</span>
                <input name="role" value="" placeholder="请输入角色，0管理员 1用户" style="height: 30px; width: 300px;"/>
            </label><br/>
            <label style="margin: 0 20px;display: block">
                <input type="submit" name="submit" onclick="relateSubmit()"
                       style="width: 80px;height: 20px; background: #ccc; border-radius: 5px;" />
            </label>
        </div>
    </div>
    <div>
        <input type="button" name="submit" onclick="tenant()" value="创建租户"
               style="width: 80px;height: 20px; background: #ccc; border-radius: 5px;color: blue;" />
        <input type="button" name="submit" onclick="user()" value="创建用户"
               style="width: 80px;height: 20px; background: #ccc; border-radius: 5px;color: blue;" />
    </div>
    <div class="tenant">
        <table border="1" width="500" align="center" cellspacing="0" cellpadding="6">
            <caption>租户列表</caption>

            <thead>
            <tr align="center">
                <th>id</th>
                <th>名称</th>
            </tr>
            </thead>

            <tbody>
            ` + tenantList + `
            </tbody>
        </table>
    </div>
    <div class="user" style="margin-top: 50px;">
        <table border="1" width="100%" align="center" cellspacing="0" cellpadding="6">
            <caption>用户列表</caption>

            <thead>
            <tr align="center">
                <th>id</th>
                <th>用户名</th>
                <th>昵称</th>
                <th>手机号</th>
                <th>性别</th>
                <th>code</th>
                <th>所属租户</th>
                <th>操作</th>
            </tr>
            </thead>

            <tbody>
            ` + userList + `
            </tbody>
        </table>
    </div>
</div>
<script>
    let userId = null;
	let host = "http://" + window.location.host;
    function relate(uid) {
        userId = uid;
        document.getElementsByClassName("pop")[0].style.display = "block";
        document.getElementsByClassName("relateBox")[0].style.display = "block";
    }

    function relateSubmit() {
        const tenantId = document.getElementsByName("tenantId")[0].value;
        if (tenantId === "") {
            document.getElementsByClassName("pop")[0].style.display = "none";
            document.getElementsByClassName("relateBox")[0].style.display = "none";
            alert('租户id不能为空');
            return
        }

        const role = document.getElementsByName("role")[0].value;
        if (role === "") {
            document.getElementsByClassName("pop")[0].style.display = "none";
            document.getElementsByClassName("relateBox")[0].style.display = "none";
            alert('角色不能为空');
            return
        }

        const data = "userId=" + userId + "&tenantId=" + tenantId + "&userType=" + role;
        createRequest(host + "/user/relate", "POST", data, function (res) {
            let result = JSON.parse(res);
            console.log(result);
            document.getElementsByClassName("pop")[0].style.display = "none";
            document.getElementsByClassName("relateBox")[0].style.display = "none";
            alert(result.message || "请求异常");
			window.location.reload();
        })
    }
    
    function tenant() {
        document.getElementsByClassName("pop")[0].style.display = "block";
        document.getElementsByClassName("tenantBox")[0].style.display = "block";
    }
    
    function tenantSubmit() {
        const tenantName = document.getElementsByName("tenantName")[0].value;
        if (tenantName === "") {
            document.getElementsByClassName("pop")[0].style.display = "none";
            document.getElementsByClassName("tenantBox")[0].style.display = "none";
            alert('租户名称不能为空');
            return
        }

        const data = "name=" + tenantName;
        createRequest(host + "/tenant/add", "POST", data, function (res) {
            let result = JSON.parse(res);
            console.log(result);
            document.getElementsByClassName("pop")[0].style.display = "none";
            document.getElementsByClassName("tenantBox")[0].style.display = "none";
            alert(result.message || "请求异常");
            window.location.reload();
        })
    }
    
    function user() {
        document.getElementsByClassName("pop")[0].style.display = "block";
        document.getElementsByClassName("userBox")[0].style.display = "block";
    }

    function userSubmit() {
        const userName = document.getElementsByName("userName")[0].value;
        if (userName === "") {
            document.getElementsByClassName("pop")[0].style.display = "none";
            document.getElementsByClassName("userBox")[0].style.display = "none";
            alert('用户名称不能为空');
            return
        }
        const nickName = document.getElementsByName("nickName")[0].value;
        if (nickName === "") {
            document.getElementsByClassName("pop")[0].style.display = "none";
            document.getElementsByClassName("userBox")[0].style.display = "none";
            alert('用户昵称不能为空');
            return
        }
        const phone = document.getElementsByName("phone")[0].value;
        if (phone === "") {
            document.getElementsByClassName("pop")[0].style.display = "none";
            document.getElementsByClassName("userBox")[0].style.display = "none";
            alert('手机号码不能为空');
            return
        }

        const data = "username=" + userName + "&nickname=" + nickName + "&phone=" + phone;
        createRequest(host + "/user/add", "POST", data, function (res) {
            let result = JSON.parse(res);
            console.log(result);
            document.getElementsByClassName("pop")[0].style.display = "none";
            document.getElementsByClassName("userBox")[0].style.display = "none";
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
        // xhr.setRequestHeader("Canhe-Control", "no-cache");//阻止浏览器读取缓存
        xhr.send(data);
    }
</script>
</body>
</html>`
}
