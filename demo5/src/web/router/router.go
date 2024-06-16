package router

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"select-course/demo5/src/constant/config"
	"select-course/demo5/src/constant/services"
	"select-course/demo5/src/web/router/middleware"
	"select-course/demo5/src/web/services/courses"
	"select-course/demo5/src/web/services/users"
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

	// tracing
	v1 := router.Group("api/v1")
	v1.Use(middleware.Auth)
	v1.Use(otelgin.Middleware(services.WebServiceName))
	v1.GET("ping/", func(context *gin.Context) {
		context.JSON(200, gin.H{
			"message": "pong",
		})
	})
	user := v1.Group("users")
	{
		// 获取用户信息
		user.GET("/", users.GetUserHandler)
	}

	coursePath := v1.Group("course")
	{
		coursePath.GET("list/", courses.GetCourseList)
		coursePath.GET("my/", courses.MyCourseList)
		coursePath.POST("select/", courses.SelectCourse)
		coursePath.POST("back/", courses.BackCourse)
	}

	return router
}
