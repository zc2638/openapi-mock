package data

/**
 * Created by zc on 2019-10-29.
 */
type ApiData struct {
	TenantId   string `json:"tenantId"`
	TenantName string `json:"tenantName"`
	UserId     string `json:"userId"`
	UserName   string `json:"userName"`
	ApiId      string `json:"apiId"`
	ApiName    string `json:"apiName"`
	ApiDesc    string `json:"apiDesc"`
	Status     int    `json:"status"`
}
