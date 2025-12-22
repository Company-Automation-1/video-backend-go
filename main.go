// Package main 应用程序入口
package main

import (
	"fmt"
	"log"

	"github.com/Company-Automation-1/video-backend-go/src/config"
	"github.com/Company-Automation-1/video-backend-go/src/infrastructure"
	"github.com/Company-Automation-1/video-backend-go/src/middleware"
	"github.com/Company-Automation-1/video-backend-go/src/query"
	"github.com/Company-Automation-1/video-backend-go/src/routes"
	"github.com/Company-Automation-1/video-backend-go/src/services"
	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)
	fmt.Println("Gin 模式:", gin.Mode())

	// 初始化 Gin 引擎
	r := gin.New()

	// 初始化数据库
	manager := infrastructure.NewDatabaseManager(cfg)
	if err := manager.InitDatabase(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}
	query.SetDefault(manager.GetDB())

	// 初始化基础设施层
	redis, err := infrastructure.NewRedis(cfg)
	if err != nil {
		log.Fatalf("初始化Redis失败: %v", err)
	}
	email := infrastructure.NewEmail(cfg)

	// 初始化业务服务层
	captchaService := services.NewCaptchaService(redis, email)
	userService := services.NewUserService(captchaService)
	authService := services.NewAuthService(&cfg.JWT, redis)

	// 注册中间件
	r.Use(middleware.CORS(&cfg.CORS))           // 跨域处理
	r.Use(middleware.ErrorRecoveryMiddleware()) // 错误恢复
	r.Use(middleware.ResponseMiddleware())      // 响应处理
	r.Use(middleware.Logger())                  // 日志记录

	// 注册路由
	pythonURL := "http://192.168.14.70:6869" // Python服务地址
	routes.RegisterRoutes(r, userService, authService, pythonURL)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		panic("服务器启动失败: " + err.Error())
	}
}
