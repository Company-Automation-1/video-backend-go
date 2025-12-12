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

// NewUserController 创建用户控制器
func NewUserController(service *services.UserService) *UserController {
	return &UserController{
		service: service,
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

// Register 用户注册
func (c *UserController) Register(ctx *gin.Context, req *dto.UserRegisterRequest) error {
	if err := c.service.Register(ctx.Request.Context(), req.Username, req.Email, req.Password, req.Captcha); err != nil {
		return err
	}
	middleware.Created(ctx, "注册成功")
	return nil
}

// SendVerificationCode 发送验证码
func (c *UserController) SendVerificationCode(ctx *gin.Context, req *dto.SendVerificationCodeRequest) error {
	if err := c.service.SendVerificationCode(ctx.Request.Context(), req.Email); err != nil {
		return err
	}
	middleware.Success(ctx, "验证码已发送")
	return nil
}

// UpdateEmail 更新邮箱
func (c *UserController) UpdateEmail(ctx *gin.Context, req *dto.UserEmailUpdateRequest) error {
	userID, err := middleware.GetUserID(ctx)
	if err != nil {
		return err
	}

	if err := c.service.UpdateEmail(ctx.Request.Context(), userID, req.Email, req.Code); err != nil {
		return err
	}
	middleware.Success(ctx, "邮箱更新成功")
	return nil
}

// Update 更新用户
func (c *UserController) Update(ctx *gin.Context, req *dto.UserUpdateRequest) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}

	user := req.ToModel()
	updatedUser, err := c.service.Update(id, user)
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
