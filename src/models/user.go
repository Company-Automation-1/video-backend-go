// Package models 定义数据模型
package models

// User 用户模型
type User struct {
	ID            uint   `gorm:"primaryKey;autoIncrement;comment:ID" json:"id"`
	Username      string `gorm:"type:varchar(100);not null;uniqueIndex:idx_username;comment:用户名" json:"username"`
	Email         string `gorm:"type:varchar(100);not null;uniqueIndex:idx_email;comment:邮箱" json:"email"`
	Password      string `gorm:"type:varchar(100);not null;comment:密码" json:"-"`
	EmailVerified bool   `gorm:"default:false;comment:邮箱是否已验证" json:"email_verified"`
	CreatedAt     int64  `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt     int64  `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "a_users"
}
