// Package services 业务逻辑服务层
package services

import (
	"github.com/Company-Automation-1/video-backend-go/src/models"
	"github.com/Company-Automation-1/video-backend-go/src/query"
	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct{}

// NewUserService 创建新的用户服务
func NewUserService() *UserService {
	return &UserService{}
}

// GetList 获取用户列表
func (s *UserService) GetList(conditions ...gen.Condition) ([]*models.User, error) {
	return query.User.Where(conditions...).Find()
}

// GetOne 获取单个用户
func (s *UserService) GetOne(conditions ...gen.Condition) (*models.User, error) {
	user, err := query.User.Where(conditions...).First()
	if err == gorm.ErrRecordNotFound {
		return nil, tools.ErrNotFound("用户不存在")
	}
	return user, err
}

// Create 创建用户
func (s *UserService) Create(user *models.User) error {
	return query.User.Create(user)
}

// Update 更新用户
func (s *UserService) Update(id uint, user *models.User) error {
	_, err := query.User.Where(query.User.ID.Eq(id)).Updates(user)
	return err
}

// Delete 删除用户
func (s *UserService) Delete(id uint) error {
	_, err := query.User.Where(query.User.ID.Eq(id)).Delete()
	return err
}
