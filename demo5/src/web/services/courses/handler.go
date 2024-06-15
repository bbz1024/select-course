package courses

import (
	"context"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"select-course/demo1/src/utils/resp"
	"select-course/demo5/src/constant/services"
	"select-course/demo5/src/rpc/course"
	"select-course/demo5/src/utils/grpc"
	logger2 "select-course/demo5/src/utils/logger"
)

var (
	courseClient course.CourseServiceClient
	logger       *zap.Logger
)

func New() {
	conn := grpc.Connect(context.Background(), services.CourseRpcServerName)
	courseClient = course.NewCourseServiceClient(conn)

	logger = logger2.LogService(services.CourseRpcServerName)
}
func GetCourseList(ctx *gin.Context) {
	courses, err := courseClient.GetAllCourses(ctx, &course.GetAllCoursesRequest{})
	if err != nil {
		panic(err)
	}
	resp.Success(ctx, courses)
}
