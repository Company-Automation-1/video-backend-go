// Package config 提供配置加载和数据库配置管理功能
package config

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	CORS     CORSConfig     `yaml:"cors"`
	Redis    RedisConfig    `yaml:"redis"`
	Email    EmailConfig    `yaml:"email"`
	JWT      JWTConfig      `yaml:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port int    `yaml:"port"`
	Mode string `yaml:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Database        string `yaml:"database"`
	Charset         string `yaml:"charset"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"`
}

// CORSConfig 跨域配置
type CORSConfig struct {
	AllowOrigins     []string `yaml:"allow_origins"`
	AllowMethods     []string `yaml:"allow_methods"`
	AllowHeaders     []string `yaml:"allow_headers"`
	ExposeHeaders    []string `yaml:"expose_headers"`
	AllowCredentials bool     `yaml:"allow_credentials"`
	MaxAge           int      `yaml:"max_age"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	SSL      bool   `yaml:"ssl"`
	Timeout  int    `yaml:"timeout"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret     string `yaml:"secret"`      // JWT密钥
	ExpireTime int    `yaml:"expire_time"` // 过期时间（小时）
}

// Load 从文件加载配置
func Load(configPath string) (*Config, error) {
	data, err := os.ReadFile(configPath) //nolint:gosec // 配置文件路径由调用方控制
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	// 设置默认值
	cfg.Server.setDefaults()
	cfg.Database.setDefaults()
	cfg.CORS.setDefaults()
	cfg.Redis.setDefaults()
	cfg.Email.setDefaults()
	cfg.JWT.setDefaults()

	return &cfg, nil
}

// setDefaults 设置服务器配置的默认值
func (c *ServerConfig) setDefaults() {
	if c.Port == 0 {
		c.Port = 8080
	}
	if c.Mode == "" {
		c.Mode = "release"
	}
}

// setDefaults 设置数据库配置的默认值
func (c *DatabaseConfig) setDefaults() {
	if c.Charset == "" {
		c.Charset = "utf8mb4"
	}
	if c.MaxIdleConns == 0 {
		c.MaxIdleConns = 10
	}
	if c.MaxOpenConns == 0 {
		c.MaxOpenConns = 100
	}
	if c.ConnMaxLifetime == 0 {
		c.ConnMaxLifetime = 3600
	}
}

// setDefaults 设置跨域配置的默认值
func (c *CORSConfig) setDefaults() {
	if c.AllowOrigins == nil {
		c.AllowOrigins = []string{"*"}
	}
	if c.AllowMethods == nil {
		c.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"}
	}
	if c.AllowHeaders == nil {
		c.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	}
	if c.ExposeHeaders == nil {
		c.ExposeHeaders = []string{"Content-Length"}
	}
	if c.MaxAge == 0 {
		c.MaxAge = 12 * 3600
	}
	// 如果 AllowOrigins 包含 "*"，强制 AllowCredentials 为 false
	if len(c.AllowOrigins) == 1 && c.AllowOrigins[0] == "*" {
		c.AllowCredentials = false
	}
}

// setDefaults 设置Redis配置的默认值
func (c *RedisConfig) setDefaults() {
	if c.Host == "" {
		c.Host = "localhost"
	}
	if c.Port == 0 {
		c.Port = 6379
	}
	if c.DB == 0 && c.DB != 0 {
		c.DB = 0
	}
}

// setDefaults 设置邮件配置的默认值
func (c *EmailConfig) setDefaults() {
	if c.Port != 465 && c.Port != 587 {
		c.Port = 587
	}
	if c.Timeout <= 0 {
		c.Timeout = 10
	}
}

// setDefaults 设置JWT配置的默认值
func (c *JWTConfig) setDefaults() {
	if c.Secret == "" {
		c.Secret = "your-secret-key-change-in-production"
	}
	if c.ExpireTime == 0 {
		c.ExpireTime = 24 // 默认24小时
	}
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	db := c.Database
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		db.Username,
		db.Password,
		db.Host,
		db.Port,
		db.Database,
		db.Charset,
	)
}
