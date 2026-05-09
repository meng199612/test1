package handlers

import (
	"strconv"

	"admin-backend/internal/database"
	"admin-backend/internal/models"
	"admin-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetPermissionsTree(c *gin.Context) {
	var perms []models.Permission
	database.DB.Order("sort ASC, id ASC").Find(&perms)

	tree := buildPermissionTree(perms, 0)
	response.Success(c, tree)
}

func buildPermissionTree(perms []models.Permission, parentID uint) []*models.Permission {
	var result []*models.Permission
	for i := range perms {
		if perms[i].ParentID == parentID {
			node := &perms[i]
			node.Children = buildPermissionTree(perms, perms[i].ID)
			result = append(result, node)
		}
	}
	return result
}

func CreatePermission(c *gin.Context) {
	var req models.Permission
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Create(&req).Error; err != nil {
		response.Error(c, "创建失败")
		return
	}

	response.SuccessWithMsg(c, "创建成功", req)
}

func UpdatePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var perm models.Permission
	if err := database.DB.First(&perm, id).Error; err != nil {
		response.Error(c, "权限不存在")
		return
	}

	if err := c.ShouldBindJSON(&perm); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Save(&perm).Error; err != nil {
		response.Error(c, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", perm)
}

func DeletePermission(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var count int64
	database.DB.Model(&models.Permission{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		response.Error(c, "存在子权限，无法删除")
		return
	}

	if err := database.DB.Delete(&models.Permission{}, id).Error; err != nil {
		response.Error(c, "删除失败")
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}
