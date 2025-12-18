// Package services 管理员服务
package services

import (
	"strings"

	"github.com/Company-Automation-1/video-backend-go/src/api/dto"
	"github.com/Company-Automation-1/video-backend-go/src/models"
	"github.com/Company-Automation-1/video-backend-go/src/query"
	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"gorm.io/gen"
	"gorm.io/gorm"
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

// GetListWithQuery 获取管理员列表（条件查询、模糊查询、范围查询）
func (s *AdminService) GetListWithQuery(queryReq *dto.AdminListQueryRequest) ([]*models.Admin, int64, error) {
	// 1. 构建基础查询条件
	conditions := tools.NewConditionBuilder().
		EqUint(&query.Admin.ID, queryReq.ID).
		EqString(&query.Admin.Username, queryReq.Username).
		Like(&query.Admin.Username, queryReq.UsernameLike).
		GteInt64(&query.Admin.CreatedAt, queryReq.CreatedAtMin).
		LteInt64(&query.Admin.CreatedAt, queryReq.CreatedAtMax).
		GteInt64(&query.Admin.UpdatedAt, queryReq.UpdatedAtMin).
		LteInt64(&query.Admin.UpdatedAt, queryReq.UpdatedAtMax).
		Build()

	q := query.Admin.Where(conditions...)

	// 2. 排序（字段存在时应用排序）
	if queryReq.OrderBy != "" {
		if orderField, ok := query.Admin.GetFieldByName(queryReq.OrderBy); ok {
			if strings.EqualFold(queryReq.Order, "asc") {
				q = q.Order(orderField.Asc())
			} else {
				q = q.Order(orderField.Desc())
			}
		}
	}

	// 3. 执行分页查询
	admins, count, err := q.FindByPage(queryReq.GetOffset(), queryReq.GetLimit())
	return admins, count, err
}

// GetOne 获取单个管理员
func (s *AdminService) GetOne(conditions ...gen.Condition) (*models.Admin, error) {
	admin, err := query.Admin.Where(conditions...).First()
	if err == gorm.ErrRecordNotFound {
		return nil, tools.ErrNotFound("管理员不存在")
	}
	return admin, err
}
