package data

/**
 * Created by zc on 2019-10-29.
 */
type ApiData struct {
	TenantName string `json:"tenantName"`
	UserName   string `json:"userName"`
	ApiId      string `json:"apiId"`
	ApiName    string `json:"apiName"`
	ApiDesc    string `json:"apiDesc"`
	Status     int    `json:"status"`
}
