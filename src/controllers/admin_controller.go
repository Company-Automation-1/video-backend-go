// Package controllers 管理员管理控制器
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

// AdminController 管理员管理控制器
type AdminController struct {
	adminService *services.AdminService
}

// NewAdminController 创建管理员管理控制器
func NewAdminController(adminService *services.AdminService) *AdminController {
	return &AdminController{
		adminService: adminService,
	}
}

// GetList 获取管理员列表（分页）
func (c *AdminController) GetList(ctx *gin.Context) error {
	// 绑定分页参数（从查询参数获取）
	var paginationReq dto.PaginationRequest
	if err := ctx.ShouldBindQuery(&paginationReq); err != nil {
		return tools.ErrBadRequest(err.Error())
	}

	// 获取分页数据
	admins, total, err := c.adminService.GetListWithPagination(
		paginationReq.GetOffset(),
		paginationReq.GetLimit(),
	)
	if err != nil {
		return err
	}

	// 转换为VO并返回分页响应
	adminList := vo.FromAdminModelList(admins)
	paginatedResp := vo.NewPaginatedResponse(
		adminList,
		paginationReq.GetPage(),
		paginationReq.GetPageSize(),
		total,
	)
	middleware.Success(ctx, paginatedResp)
	return nil
}

// GetProfile 获取个人信息（当前登录管理员）
func (c *AdminController) GetProfile(ctx *gin.Context) error {
	adminID, err := middleware.GetAdminID(ctx)
	if err != nil {
		return err
	}
	admin, err := c.adminService.GetOne(query.Admin.ID.Eq(adminID))
	if err != nil {
		return err
	}
	middleware.Success(ctx, vo.FromAdminModel(admin))
	return nil
}
