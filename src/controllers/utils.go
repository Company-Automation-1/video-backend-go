// Package controllers 控制器工具函数
package controllers

import (
	"strconv"

	"github.com/Company-Automation-1/video-backend-go/src/tools"
	"github.com/gin-gonic/gin"
)

// parseID 解析路径参数中的ID
func parseID(ctx *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 32)
	if err != nil {
		return 0, tools.ErrBadRequest("无效的用户ID")
	}
	return uint(id), nil
}
