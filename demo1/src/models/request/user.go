package request

type UserReq struct {
	UserID int `json:"user_id" form:"user_id" binding:"required"`
}
