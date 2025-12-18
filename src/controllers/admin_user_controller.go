// Package controllers 管理员用户管理控制器
package controllers

import (
	"github.com/Company-Automation-1/video-backend-go/src/api/dto"
	"github.com/Company-Automation-1/video-backend-go/src/api/vo"
	"github.com/Company-Automation-1/video-backend-go/src/middleware"
	"github.com/Company-Automation-1/video-backend-go/src/query"
	"github.com/Company-Automation-1/video-backend-go/src/services"
	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"github.com/gin-gonic/gin"
)

// AdminUserController 管理员用户管理控制器
type AdminUserController struct {
	userService *services.UserService
}

// NewAdminUserController 创建管理员用户管理控制器
func NewAdminUserController(userService *services.UserService) *AdminUserController {
	return &AdminUserController{
		userService: userService,
	}
}

// GetList 获取用户列表（管理员权限，分页，支持条件查询、模糊查询、范围查询）
func (c *AdminUserController) GetList(ctx *gin.Context) error {
	// 绑定查询参数（从查询参数获取）
	var queryReq dto.UserListQueryRequest
	if err := ctx.ShouldBindQuery(&queryReq); err != nil {
		return tools.ErrBadRequest(err.Error())
	}

	// 调用服务层获取数据
	users, total, err := c.userService.GetListWithQuery(&queryReq)
	if err != nil {
		return err
	}

	// 转换为VO并返回分页响应
	userList := vo.FromModelList(users)
	paginatedResp := vo.NewPaginatedResponse(
		userList,
		queryReq.GetPage(),
		queryReq.GetPageSize(),
		total,
	)
	middleware.Success(ctx, paginatedResp)
	return nil
}

// GetOne 获取单个用户详情（管理员权限）
func (c *AdminUserController) GetOne(ctx *gin.Context) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}
	user, err := c.userService.GetOne(query.User.ID.Eq(id))
	if err != nil {
		return err
	}
	middleware.Success(ctx, vo.FromModel(user))
	return nil
}

// Update 更新用户信息（包括积分）（管理员权限）
func (c *AdminUserController) Update(ctx *gin.Context, req *dto.UserUpdateRequest) error {
	id, err := parseID(ctx)
	if err != nil {
		return err
	}

	user := req.ToModel()
	updatedUser, err := c.userService.Update(id, user)
	if err != nil {
		return err
	}

	middleware.Success(ctx, vo.FromModel(updatedUser))
	return nil
}
