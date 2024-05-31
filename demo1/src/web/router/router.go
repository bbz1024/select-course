package router

import (
	"github.com/gin-gonic/gin"
	"select-course/demo1/src/constant/config"
	"select-course/demo1/src/web/router/middleware"
	"select-course/demo1/src/web/services/course"
	"select-course/demo1/src/web/services/users"
)

func InitApiRouter() *gin.Engine {
	var router *gin.Engine
	if config.EnvCfg.ProjectMode == "dev" {
		gin.SetMode(gin.DebugMode)
		gin.ForceConsoleColor()
		router = gin.Default()

	} else {
		gin.SetMode(gin.ReleaseMode)
		router = gin.New()
		router.Use(gin.Logger())
	}
	v1 := router.Group("api/v1")
	v1.Use(middleware.Auth)
	user := v1.Group("user")

	{

		// 获取用户信息
		user.GET("/", users.GetUserHandler)
	}

	coursePath := v1.Group("course")
	{
		coursePath.GET("list/", course.GetCourseList)
		//course.GET("my/", nil)
		//course.POST("select/", nil)
	}

	return router
}
