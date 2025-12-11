// Package controllers 用户控制器
package controllers

import (
	"github.com/Company-Automation-1/video-backend-go/src/api/dto"
	"github.com/Company-Automation-1/video-backend-go/src/api/vo"
	"github.com/Company-Automation-1/video-backend-go/src/middleware"
	"github.com/Company-Automation-1/video-backend-go/src/query"
	"github.com/Company-Automation-1/video-backend-go/src/services"
	"github.com/gin-gonic/gin"
)

// UserController 用户控制器
type UserController struct {
	service *services.UserService
}

// NewUserController 创建新的用户控制器
func NewUserController() *UserController {
	return &UserController{
		service: services.NewUserService(),
	}
}

// GetList 获取用户列表
func (c *UserController) GetList(ctx *gin.Context) error {
	users, err := c.service.GetList()
	if err != nil {
		return err
	}
	middleware.Success(ctx, vo.FromModelList(users))
	return nil
}

// GetOne 获取单个用户
func (c *UserController) GetOne(ctx *gin.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}
	user, err := c.service.GetOne(query.User.ID.Eq(id))
	if err != nil {
		return err
	}
	middleware.Success(ctx, vo.FromModel(user))
	return nil
}

// Create 创建用户
func (c *UserController) Create(ctx *gin.Context, req *dto.UserCreateRequest) error {
	user := req.ToModel()
	if err := c.service.Create(user); err != nil {
		return err
	}
	middleware.Created(ctx, vo.FromModel(user))
	return nil
}

// Update 更新用户
func (c *UserController) Update(ctx *gin.Context, req *dto.UserUpdateRequest) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}

	user := req.ToModel()
	if err := c.service.Update(id, user); err != nil {
		return err
	}

	updatedUser, err := c.service.GetOne(query.User.ID.Eq(id))
	if err != nil {
		return err
	}
	middleware.Success(ctx, vo.FromModel(updatedUser))
	return nil
}

// Delete 删除用户
func (c *UserController) Delete(ctx *gin.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}
	if err := c.service.Delete(id); err != nil {
		return err
	}
	middleware.Success(ctx, "删除成功")
	return nil
}
