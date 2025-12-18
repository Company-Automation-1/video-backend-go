// Package vo 管理员值对象
package vo

import (
	"github.com/Company-Automation-1/video-backend-go/src/models"
)

// AdminVO 管理员值对象
type AdminVO struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

// FromAdminModel 从模型转换为VO
func FromAdminModel(admin *models.Admin) *AdminVO {
	return &AdminVO{
		ID:        admin.ID,
		Username:  admin.Username,
		CreatedAt: admin.CreatedAt,
		UpdatedAt: admin.UpdatedAt,
	}
}

// FromAdminModelList 从模型列表转换为VO列表
func FromAdminModelList(admins []*models.Admin) []*AdminVO {
	result := make([]*AdminVO, len(admins))
	for i, admin := range admins {
		result[i] = FromAdminModel(admin)
	}
	return result
}
