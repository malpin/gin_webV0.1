package Bean

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Result struct {
	Code int64       `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

// Success 定义统一的成功返回函数
func Success(c *gin.Context, data interface{}) {
	result := newDataResult(SUCCESS, data)
	c.JSON(http.StatusOK, result)
}

// Error 统一错误返回方法，所有错误都调用此方法
func Error(c *gin.Context, code CodeMsg) {
	result := newErrResult(code)
	c.JSON(http.StatusOK, result)
}

// ErrorWithMsg 统一错误返回方法，所有自定义的错误都调用此方法
func ErrorWithMsg(c *gin.Context, code CodeMsg, msg interface{}) {
	result := newErrWithMsgResult(code, msg)
	c.JSON(http.StatusOK, result)
}

// newDataResult 使用预定义的状态码和指定的数据信息 返回一个*Result
func newDataResult(code CodeMsg, data interface{}) *Result {
	return &Result{
		Code: code.code,
		Msg:  code.msg,
		Data: data,
	}
}

// newErrResult 使用预定义的状态码和错误信息 返回一个*Result
func newErrResult(code CodeMsg) *Result {
	return &Result{
		Code: code.code,
		Msg:  code.msg,
	}
}

// newErrWithMsgResult 使用预定义的状态码和自定义的错误信息 返回一个*Result
func newErrWithMsgResult(code CodeMsg, msg interface{}) *Result {
	return &Result{
		Code: code.code,
		Msg:  msg,
	}
}

/*


type Result struct {
	Code ResMsg      `json:"code"`
	Msg  interface{} `json:"Msg"`
	Data interface{} `json:"Data"`
}

// Success 定义统一的成功返回函数
func Success(c *gin.Context, data interface{}) {
	result := newDataResult(data, SUCCESS)
	c.JSON(http.StatusOK, result)
}

// Error 统一错误返回方法，所有错误都调用此方法
func Error(c *gin.Context, code ResMsg) {
	result := newErrResult(code)
	c.JSON(http.StatusOK, result)
}

// ErrorWithMsg 统一错误返回方法，所有自定义的错误都调用此方法
func ErrorWithMsg(c *gin.Context, code ResMsg, msg interface{}) {
	result := newErrWithMsgResult(code, msg)
	c.JSON(http.StatusOK, result)
}

func newDataResult(data interface{}, code ResMsg) *Result {
	return &Result{
		Code: code,
		Msg:  code.GetMsg(),
		Data: data,
	}
}

func newErrResult(code ResMsg) *Result {
	return &Result{
		Code: code,
		Msg:  code.GetMsg(),
	}
}

func newErrWithMsgResult(code ResMsg, msg interface{}) *Result {
	return &Result{
		Code: code,
		Msg:  msg,
	}
}

*/
