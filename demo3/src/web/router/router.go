package router

import (
	"github.com/gin-gonic/gin"
	"select-course/demo3/src/constant/config"
	"select-course/demo3/src/web/router/middleware"
	"select-course/demo3/src/web/services/course"
	"select-course/demo3/src/web/services/users"
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
	v1.GET("ping/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	user := v1.Group("user")

	{

		// 获取用户信息
		user.GET("/", users.GetUserHandler)
	}

	coursePath := v1.Group("course")
	{
		coursePath.GET("list/", course.GetCourseList)
		coursePath.GET("my/", course.MyCourseList)
		coursePath.POST("select/", course.SelectCourse)
		coursePath.POST("back/", course.BackCourse)
	}

	return router
}
