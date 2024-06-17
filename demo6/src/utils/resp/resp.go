package resp

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"net/http"
	"select-course/demo6/src/constant/code"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	StatusMsg  string `json:"status_msg"`
}

func Fail(ctx *gin.Context, statusCode int, StatusMsg string) {
	ctx.JSON(
		http.StatusOK, &Response{
			statusCode, StatusMsg,
		},
	)
}
func ParamErr(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK, &Response{
			code.ParamErr, code.ParamErrMsg,
		},
	)
}
func DBError(ctx *gin.Context) {
	ctx.JSON(http.StatusInternalServerError, &Response{
		code.DBError, code.DBErrorMsg,
	})
}

type CustomJSON struct {
	Data    proto.Message
	Context *gin.Context
}

var m = protojson.MarshalOptions{
	EmitUnpopulated: true,
	UseProtoNames:   true,
}

func (r CustomJSON) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	res, _ := m.Marshal(r.Data)
	_, err = w.Write(res)
	return
}

func (r CustomJSON) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
}
