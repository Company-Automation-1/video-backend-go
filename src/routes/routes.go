// Package routes 路由注册
package routes

import (
	"github.com/Company-Automation-1/video-backend-go/src/controllers"
	"github.com/Company-Automation-1/video-backend-go/src/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine) {
	// API v1 路由组
	v1 := r.Group("/api/v1")
	// 用户路由
	userController := controllers.NewUserController()
	users := v1.Group("/users")
	users.GET("", middleware.Handle(userController.GetList))
	users.GET("/:id", middleware.Handle(userController.GetOne))
	users.POST("", middleware.Bind(userController.Create))
	users.PUT("/:id", middleware.Bind(userController.Update))
	users.DELETE("/:id", middleware.Handle(userController.Delete))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		middleware.Success(c, "ok")
	})
}
