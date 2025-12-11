// Package main 应用程序入口
package main

import (
	"fmt"
	"log"

	"github.com/Company-Automation-1/video-backend-go/src/config"
	"github.com/Company-Automation-1/video-backend-go/src/database"
	"github.com/Company-Automation-1/video-backend-go/src/middleware"
	"github.com/Company-Automation-1/video-backend-go/src/query"
	"github.com/Company-Automation-1/video-backend-go/src/routes"
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

	// 初始化 Gin 引擎（不使用 Default，避免重复中间件）
	r := gin.New()

	// 初始化配置管理器
	manager := database.NewManager(cfg)

	// 初始化数据库连接
	if err := manager.InitDatabase(); err != nil {
		log.Fatalf("初始化数据库失败: %v", err)
	}

	// 初始化 GORM Gen 查询
	query.SetDefault(manager.GetDB())

	// 中间件
	r.Use(middleware.CORS(&cfg.CORS))           // CORS 中间件（应该放在最前面，优先处理 OPTIONS 请求）
	r.Use(middleware.ErrorRecoveryMiddleware()) // 全局错误恢复中间件
	r.Use(middleware.ResponseMiddleware())      // 响应包装中间件
	r.Use(middleware.Logger())                  // Logger 中间件（所有模式都显示请求日志）

	// 注册路由
	routes.RegisterRoutes(r)

	// 启动服务器
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
