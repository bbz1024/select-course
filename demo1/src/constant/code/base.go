package code

const (
	Success  int = 200
	ParamErr     = 400
	NotFound     = 404
	Fail         = 500

	// user

	UserNotFound = 1000

	// mysql

	DBError = 2000
)
const (
	SuccessMsg      string = "success"
	ParamErrMsg            = "param error"
	NotFoundMsg            = "not found"
	FailMsg                = "fail"
	UserNotFoundMsg        = "用户不存在"
	DBErrorMsg             = "查询失败"
)
