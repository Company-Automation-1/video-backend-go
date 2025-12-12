// Package infrastructure 基础设施层：数据库管理
package infrastructure

import (
	"time"

	"github.com/Company-Automation-1/video-backend-go/src/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Manager 配置和数据库管理器
type Manager struct {
	cfg *config.Config
	db  *gorm.DB
}

// NewDatabaseManager 创建新的数据库管理器
func NewDatabaseManager(cfg *config.Config) *Manager {
	return &Manager{
		cfg: cfg,
	}
}

// InitDatabase 初始化数据库连接
func (m *Manager) InitDatabase() error {
	dsn := m.cfg.GetDSN()

	// 根据服务器模式设置 GORM 日志级别
	var logLevel logger.LogLevel
	switch m.cfg.Server.Mode {
	case "debug":
		logLevel = logger.Info // Debug 模式：显示所有 SQL 查询和错误
	default:
		logLevel = logger.Silent // Release 模式：完全静默（业务层已处理错误）
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(m.cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(m.cfg.Database.ConnMaxLifetime) * time.Second)

	m.db = db
	return nil
}

// GetDB 获取数据库连接
func (m *Manager) GetDB() *gorm.DB {
	return m.db
}

// GetConfig 获取配置
func (m *Manager) GetConfig() *config.Config {
	return m.cfg
}
