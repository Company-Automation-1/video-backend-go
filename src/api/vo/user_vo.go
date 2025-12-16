// Package vo 值对象，负责出参格式化
package vo

import (
	"github.com/Company-Automation-1/video-backend-go/src/models"
)

// UserVO 用户值对象
type UserVO struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Points    int    `json:"points"` // 前端显示：NULL 显示为 0
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// FromModel 从模型转换为VO
func FromModel(user *models.User) *UserVO {
	points := 0
	if user.Points != nil {
		points = *user.Points
	}
	return &UserVO{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Points:    points, // NULL 时显示为 0
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

// FromModelList 从模型列表转换为VO列表
func FromModelList(users []*models.User) []*UserVO {
	result := make([]*UserVO, len(users))
	for i, user := range users {
		result[i] = FromModel(user)
	}
	return result
}
