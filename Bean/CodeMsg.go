package Bean

import "errors"

type CodeMsg struct {
	code      int64
	msg       interface{}
	MarkError error
}

// 通用错误码定义
// 处理成功消息码
var SUCCESS = CodeMsg{200, "success", nil}

//非法数据错误码
var DATA_ERROR = CodeMsg{1000, "非法数据", errors.New("非法的数据")}
var SYSTEM_BUSY = CodeMsg{1001, "系统繁忙", errors.New("系统繁忙")}
var CPACHA_EMPTY = CodeMsg{1002, "验证码不能为空", errors.New("验证码不能为空")}
var CPACHA_ERROR = CodeMsg{1003, "验证码错误", errors.New("验证码错误")}
var SESSION_EXPIRED = CodeMsg{1004, "会话已过期，请重新登录", errors.New("会话已过期，请刷新页面重试")}
var USER_SESSION_EXPIRED = CodeMsg{1005, "还未登录或会话失效，请重新登录！", errors.New("还未登录或会话失效，请重新登录！")}
var COMMUNITY_ID_ERROR = CodeMsg{1006, "社区的id错误", errors.New("社区的id错误")}

//用户管理类错误
var USERNAME_EMPTY = CodeMsg{2000, "用户名不能为空", errors.New("用户名不能为空")}
var PASSWORD_EMPTY = CodeMsg{2001, "密码不能为空", errors.New("密码不能为空")}
var USERNAME_EXIST = CodeMsg{2002, "用户已经存在了", errors.New("用户已经存在了")}
var USERNAME_NO_EXIST = CodeMsg{2003, "该用户名不存在", errors.New("该用户名不存在")}
var PASSWORD_ERROR = CodeMsg{2004, "用户名或密码错误", errors.New("用户名或密码错误")}
