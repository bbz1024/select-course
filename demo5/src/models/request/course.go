package request

type SelectCourseReq struct {
	UserID   uint `json:"user_id" form:"user_id" binding:"required"`
	CourseID uint `json:"course_id" form:"course_id" binding:"required"`
}
