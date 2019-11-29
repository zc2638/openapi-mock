package service

/**
 * Created by zc on 2019-10-24.
 */
type BaseService struct{}

type RequestError string

func (e RequestError) Error() string { return string(e) }

const (
	ErrorRequest = RequestError("请求异常")
	AuthCodeError = RequestError("authCode错误")
	TokenError = RequestError("token异常")
	AuthError = RequestError("身份认证失败")
	TenantError = RequestError("租户认证失败")
	TenantRepeat = RequestError("租户已存在")
	TenantNotExist = RequestError("租户不存在")
	UserRepeat = RequestError("用户名已存在")
	UserNotExist = RequestError("用户不存在")
	UserRelateRepeat = RequestError("用户租户已关联")
	PhoneRepeat = RequestError("联系方式已存在")
	APIExist = RequestError("api发起申请或已上架")
	APINotExist = RequestError("api不存在")
)