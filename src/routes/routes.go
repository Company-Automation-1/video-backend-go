// Package routes 路由注册
package routes

import (
	"github.com/Company-Automation-1/video-backend-go/src/controllers"
	"github.com/Company-Automation-1/video-backend-go/src/middleware"
	"github.com/Company-Automation-1/video-backend-go/src/services"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(
	r *gin.Engine,
	userService *services.UserService,
	authService *services.AuthService,
) {
	// API v1 路由组
	v1 := r.Group("/api/v1")

	// 认证路由（无需鉴权）
	authController := controllers.NewAuthController(authService)
	auth := v1.Group("/auth")
	auth.POST("/user/login", middleware.Bind(authController.UserLogin))
	auth.POST("/admin/login", middleware.Bind(authController.AdminLogin))
	auth.POST("/logout", middleware.AuthMiddleware(authService), middleware.Handle(authController.Logout))

	// 用户路由
	userController := controllers.NewUserController(userService)
	users := v1.Group("/users")

	// 公开路由
	users.POST("/send-verification-code", middleware.Bind(userController.SendVerificationCode))
	users.POST("/register", middleware.Bind(userController.Register))

	// 需要本人权限的路由
	users.GET("/profile", middleware.AuthMiddleware(authService), middleware.Handle(userController.GetProfile))
	users.POST("/update-email", middleware.AuthMiddleware(authService), middleware.Bind(userController.UpdateEmail))
	users.PUT("/:id", middleware.SelfMiddleware(authService), middleware.Bind(userController.Update))
	users.DELETE("/:id", middleware.SelfMiddleware(authService), middleware.Handle(userController.Delete))

	// 管理员路由
	admin := v1.Group("/admin")
	
	// 管理员管理路由（需要管理员认证）
	adminService := services.NewAdminService()
	adminController := controllers.NewAdminController(adminService)
	admin.GET("/profile", middleware.AdminMiddleware(authService), middleware.Handle(adminController.GetProfile))
	admins := admin.Group("/admins")
	admins.Use(middleware.AdminMiddleware(authService))
	admins.GET("", middleware.Handle(adminController.GetList))
	
	// 管理员用户管理路由（需要管理员认证）
	adminUserController := controllers.NewAdminUserController(userService)
	adminUsers := admin.Group("/users")
	adminUsers.Use(middleware.AdminMiddleware(authService))
	adminUsers.GET("", middleware.Handle(adminUserController.GetList))
	adminUsers.GET("/:id", middleware.Handle(adminUserController.GetOne))
	adminUsers.PUT("/:id", middleware.Bind(adminUserController.Update))

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		middleware.Success(c, "ok")
	})
}
