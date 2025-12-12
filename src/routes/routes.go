// Package routes 路由注册
package routes

import (
	"github.com/Company-Automation-1/video-backend-go/src/controllers"
	"github.com/Company-Automation-1/video-backend-go/src/middleware"
	"github.com/Company-Automation-1/video-backend-go/src/services"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, userService *services.UserService, authService *services.AuthService) {
	// API v1 路由组
	v1 := r.Group("/api/v1")

	// 认证路由（无需鉴权）
	authController := controllers.NewAuthController(authService)
	auth := v1.Group("/auth")
	auth.POST("/login", middleware.Bind(authController.Login))
	auth.POST("/logout", middleware.AuthMiddleware(authService), middleware.Handle(authController.Logout))

	// 用户路由
	userController := controllers.NewUserController(userService)
	users := v1.Group("/users")

	// 公开路由
	users.POST("/send-verification-code", middleware.Bind(userController.SendVerificationCode))
	users.POST("/register", middleware.Bind(userController.Register))

	// 需要鉴权的路由
	users.Use(middleware.AuthMiddleware(authService))
	users.GET("", middleware.Handle(userController.GetList))
	users.GET("/:id", middleware.Handle(userController.GetOne))
	users.POST("/update-email", middleware.Bind(userController.UpdateEmail))
	users.PUT("/:id", middleware.Bind(userController.Update))
	users.DELETE("/:id", middleware.Handle(userController.Delete))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		middleware.Success(c, "ok")
	})
}
