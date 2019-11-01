package template

/**
 * Created by zc on 2019-11-01.
 */
func Error(err error) string {
	return `<!doctype html>
<html xmlns=http://www.w3.org/1999/xhtml>
<meta charset=utf-8>
<title>OpenAPI-mock用户中心</title>
<body>
<div class="content">
	<p>出错啦，哈哈~</p>
	<p>` + err.Error() + `</p>
</div>
</body>
</html>`
}