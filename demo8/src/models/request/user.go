package request

type UserReq struct {
	UserID int64 `json:"user_id" form:"user_id" binding:"required"`
}
