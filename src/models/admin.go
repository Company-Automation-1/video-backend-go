// Package models 定义数据模型
package models

// Admin 管理员模型
type Admin struct {
	ID        uint   `gorm:"primaryKey;autoIncrement;comment:ID" json:"id"`
	Username  string `gorm:"type:varchar(100);not null;uniqueIndex:idx_username;comment:用户名" json:"username"`
	Password  string `gorm:"type:varchar(100);not null;comment:密码" json:"-"`
	CreatedAt int64  `gorm:"autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt int64  `gorm:"autoUpdateTime;comment:更新时间" json:"updated_at"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "a_admins"
}

