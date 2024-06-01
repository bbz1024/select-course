package request

type SelectCourseReq struct {
	UserID   int `json:"user_id" form:"user_id" binding:"required"`
	CourseID int `json:"course_id" form:"course_id" binding:"required"`
}
