package template

import "mock/data"

/**
 * Created by zc on 2019-11-25.
 */
func ApiListTemplate(apiSet []data.ApiData) string {

	apiList := ""
	for _, api := range apiSet {
		status := "待审核"
		submit := `<td>
                    <button onclick="auditSubmit('`+ api.ApiId +`', '1')">同意</button>
                    <button onclick="auditSubmit('`+ api.ApiId +`', '2')">拒绝</button>
                </td>`
		switch api.Status {
		case 1:
			status = "已上架"
			submit = `<td>
                    <button onclick="auditSubmit('`+ api.ApiId +`', '3')">下架</button>
                </td>`
		case 2:
			status = "已拒绝"
		}



		apiList += `<tr align="center">
                <td>` + api.ApiId + `</td>
                <td>` + api.ApiName + `</td>
                <td>` + api.TenantName + `</td>
                <td>` + api.UserName + `</td>
                <td>` + api.ApiDesc + `</td>
                <td>` + status + `</td>
                ` + submit + `
            </tr>`
	}

	return `<!doctype html>
<html xmlns=http://www.w3.org/1999/xhtml>
<meta charset=utf-8>
<title>OpenAPI-mock Store</title>
<body>
<div class="content">
	<div>
        <input type="button" name="submit" onclick="buy()" value="购买API"
               style="width: 80px;height: 20px; background: #ccc; border-radius: 5px;color: blue;cursor: pointer;" />
    </div>
    <div class="api">
        <table border="1" width="100%" align="center" cellspacing="0" cellpadding="6">
            <caption>API列表</caption>

            <thead>
            <tr align="center">
                <th>id</th>
                <th>名称</th>
                <th>租户名</th>
                <th>用户名</th>
                <th>描述</th>
                <th>状态</th>
                <th>操作</th>
            </tr>
            </thead>

            <tbody>
            ` + apiList + `
            </tbody>
        </table>
    </div>
</div>
<script>
    let host = "http://" + window.location.host;
    function auditSubmit(apiId, status) {
        const data = "apiId=" + apiId + "&status=" + status;
        createRequest(host + "/mock/api/audit", "POST", data, function (res) {
            let result = JSON.parse(res);
            console.log(result);
            alert(result.message || "请求异常");
            window.location.reload();
        })
    }

	function buy() {
		window.location.href = host + "/mock/contract";
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