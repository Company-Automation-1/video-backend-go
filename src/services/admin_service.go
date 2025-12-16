// Package services 管理员服务
package services

import (
	"github.com/Company-Automation-1/video-backend-go/src/models"
	"github.com/Company-Automation-1/video-backend-go/src/query"
	"gorm.io/gen"
)

// AdminService 管理员服务
type AdminService struct {
}

// NewAdminService 创建管理员服务
func NewAdminService() *AdminService {
	return &AdminService{}
}

// GetListWithPagination 获取管理员列表（分页）
func (s *AdminService) GetListWithPagination(
	offset, limit int,
	conditions ...gen.Condition,
) ([]*models.Admin, int64, error) {
	admins, count, err := query.Admin.Where(conditions...).FindByPage(offset, limit)
	return admins, count, err
}
