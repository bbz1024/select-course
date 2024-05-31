package course

import (
	"github.com/gin-gonic/gin"
	"select-course/demo1/src/models"
	"select-course/demo1/src/storage/database"
	"select-course/demo1/src/utils/resp"
)

func GetCourseList(ctx *gin.Context) {
	var courseList []*models.Course
	if err := database.Client.Find(&courseList).Error; err != nil {
		resp.DBError(ctx)
		return
	}
	resp.Success(ctx, courseList)
}
