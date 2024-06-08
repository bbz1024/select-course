package resp

import (
	"github.com/gin-gonic/gin"
	"select-course/demo4/src/constant/code"
)

type Response struct {
	Code  int         `json:"code"`
	Data  interface{} `json:"data"`
	Msg   string      `json:"msg"`
	Error string      `json:"error"`
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(code.Success, &Response{
		Code: code.Success,
		Data: data,
		Msg:  code.SuccessMsg,
	})

}

func Fail(ctx *gin.Context, status int, cd int, err string) {
	ctx.JSON(
		status, &Response{
			Code:  cd,
			Error: err,
			Data:  nil,
		},
	)
}
func ParamErr(ctx *gin.Context) {
	Fail(ctx, code.ParamErr, code.ParamErr, code.ParamErrMsg)
}
func DBError(ctx *gin.Context) {
	Fail(ctx, code.Fail, code.DBError, code.DBErrorMsg)
}
