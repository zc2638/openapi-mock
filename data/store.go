package data

/**
 * Created by zc on 2019-10-29.
 */
type ApiData struct {
	TenantId   string `json:"tenantId" vdr:"required;msg=缺少租户标识"`
	TenantName string `json:"tenantName" vdr:"required;msg=缺少租户名称"`
	UserId     string `json:"userId" vdr:"required;msg=缺少用户标识"`
	UserName   string `json:"userName" vdr:"required;msg=缺少用户名称"`
	ApiId      string `json:"apiId" vdr:"required;msg=缺少api id参数"`
	ApiName    string `json:"apiName" vdr:"required;msg=缺少api名称"`
	ApiDesc    string `json:"apiDesc" vdr:"required;msg=缺少api描述"`
	ApiLogo    string `json:"apiLogo"`
	Status     int    `json:"status"`
}

type Contract struct {
	Email            string     `json:"email"`
	Username         string     `json:"username"`
	Phone            string     `json:"phone"`
	Gender           int        `json:"gender"`
	UucUserId        string     `json:"uucUserId"`
	CubaTenantId     string     `json:"cubaTenantId"`
	CubaTenantName   string     `json:"cubaTenantName"`
	CubaUserType     int        `json:"cubaUserType"`
	ProviderTenantId string     `json:"providerTenantId"`
	ApiOfSCs         []ApiOfSCs `json:"apiOfSCs"`
}

type ApiOfSCs struct {
	IdApi             string `json:"idApi"`
	TotalTimes        int    `json:"totalTimes"`
	TenantEverySecond int    `json:"tenantEverySecond"`
	TenantEveryMinute int    `json:"tenantEveryMinute"`
	TenantEveryHour   int    `json:"tenantEveryHour"`
	TenantEveryDay    int    `json:"tenantEveryDay"`
}
