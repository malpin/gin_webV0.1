package Bean

type ResMsg int64

const (
	SUCCESS                 ResMsg = 0
	DATA_ERROR              ResMsg = -1
	SYSTEM_BUSY             ResMsg = -2
	CPACHA_EMPTY            ResMsg = -3
	SESSION_EXPIRED         ResMsg = -4
	CPACHA_ERROR            ResMsg = -5
	USER_SESSION_EXPIRED    ResMsg = -6
	ADMIN_USERNAME_EMPTY    ResMsg = -7
	ADMIN_PASSWORD_EMPTY    ResMsg = -8
	ADMIN_USERNAME_NO_EXIST ResMsg = -9
	ADMIN_USERNAME_EXIST    ResMsg = -10
	ADMIN_PASSWORD_ERROR    ResMsg = -11
)

var codem = map[ResMsg]string{
	//通用错误码定义
	//处理成功消息码
	SUCCESS: "success",
	//非法数据错误码
	DATA_ERROR:           "非法数据",
	SYSTEM_BUSY:          "系统繁忙",
	CPACHA_EMPTY:         "验证码不能为空！",
	CPACHA_ERROR:         "验证码错误！",
	SESSION_EXPIRED:      "会话已失效，请刷新页面重试！",
	USER_SESSION_EXPIRED: "还未登录或会话失效，请重新登录！",
	//后台管理类错误码
	//用户管理类错误
	ADMIN_USERNAME_EMPTY: "用户名不能为空！",
	ADMIN_PASSWORD_EMPTY: "密码不能为空！",
	//登录类错误码
	ADMIN_USERNAME_EXIST:    "该用户名已存在！",
	ADMIN_USERNAME_NO_EXIST: "该用户名不存在！",
	ADMIN_PASSWORD_ERROR:    "密码错误！",
}

func (c ResMsg) GetMsg() string {
	return codem[c]
}
