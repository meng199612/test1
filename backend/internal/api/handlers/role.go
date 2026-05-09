package handlers

import (
	"strconv"

	"admin-backend/internal/database"
	"admin-backend/internal/models"
	"admin-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetRoles(c *gin.Context) {
	var roles []models.Role
	database.DB.Order("id ASC").Find(&roles)
	response.Success(c, roles)
}

func GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var role models.Role
	if err := database.DB.Preload("Permissions").Preload("Menus").First(&role, id).Error; err != nil {
		response.Error(c, "角色不存在")
		return
	}

	response.Success(c, role)
}

func CreateRole(c *gin.Context) {
	var req models.Role
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var existing models.Role
	if database.DB.Where("code = ?", req.Code).First(&existing).Error == nil {
		response.Error(c, "角色编码已存在")
		return
	}

	if err := database.DB.Create(&req).Error; err != nil {
		response.Error(c, "创建失败")
		return
	}

	response.SuccessWithMsg(c, "创建成功", req)
}

func UpdateRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		response.Error(c, "角色不存在")
		return
	}

	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Description string `json:"description"`
		Status      int8   `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	updates["status"] = req.Status

	if err := database.DB.Model(&role).Updates(updates).Error; err != nil {
		response.Error(c, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", role)
}

func DeleteRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		response.Error(c, "角色不存在")
		return
	}

	var count int64
	database.DB.Table("user_roles").Where("role_id = ?", id).Count(&count)
	if count > 0 {
		response.Error(c, "该角色已被用户使用，无法删除")
		return
	}

	if err := database.DB.Delete(&role).Error; err != nil {
		response.Error(c, "删除失败")
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

func AssignPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		response.Error(c, "角色不存在")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var perms []models.Permission
	if len(req.PermissionIDs) > 0 {
		database.DB.Where("id IN ?", req.PermissionIDs).Find(&perms)
	}

	database.DB.Model(&role).Association("Permissions").Replace(&perms)

	response.SuccessWithMsg(c, "分配成功", nil)
}

func AssignMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var role models.Role
	if err := database.DB.First(&role, id).Error; err != nil {
		response.Error(c, "角色不存在")
		return
	}

	var req struct {
		MenuIDs []uint `json:"menu_ids"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var menus []models.Menu
	if len(req.MenuIDs) > 0 {
		database.DB.Where("id IN ?", req.MenuIDs).Find(&menus)
	}

	database.DB.Model(&role).Association("Menus").Replace(&menus)

	response.SuccessWithMsg(c, "分配成功", nil)
}
