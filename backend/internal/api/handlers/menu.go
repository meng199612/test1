package handlers

import (
	"strconv"

	"admin-backend/internal/database"
	"admin-backend/internal/models"
	"admin-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetMenusTree(c *gin.Context) {
	var menus []models.Menu
	database.DB.Where("visible = ?", 1).Order("sort ASC, id ASC").Find(&menus)

	tree := buildMenuTree(menus, 0)
	response.Success(c, tree)
}

func GetAllMenus(c *gin.Context) {
	var menus []models.Menu
	database.DB.Order("sort ASC, id ASC").Find(&menus)

	tree := buildMenuTree(menus, 0)
	response.Success(c, tree)
}

func GetUserMenus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	var roles []models.Role
	database.DB.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Preload("Menus", "visible = ?", 1).
		Find(&roles)

	menuMap := make(map[uint]*models.Menu)
	for _, role := range roles {
		for i := range role.Menus {
			menu := &role.Menus[i]
			menuMap[menu.ID] = menu
		}
	}

	var menus []models.Menu
	for _, m := range menuMap {
		menus = append(menus, *m)
	}

	if len(menuMap) == 0 {
		database.DB.Where("visible = ?", 1).Order("sort ASC, id ASC").Find(&menus)
	}

	tree := buildMenuTree(menus, 0)
	response.Success(c, tree)
}

func buildMenuTree(menus []models.Menu, parentID uint) []*models.Menu {
	var result []*models.Menu
	for i := range menus {
		if menus[i].ParentID == parentID {
			node := &menus[i]
			node.Children = buildMenuTree(menus, menus[i].ID)
			result = append(result, node)
		}
	}
	return result
}

func GetMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var menu models.Menu
	if err := database.DB.First(&menu, id).Error; err != nil {
		response.Error(c, "菜单不存在")
		return
	}

	response.Success(c, menu)
}

func CreateMenu(c *gin.Context) {
	var req models.Menu
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

func UpdateMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var menu models.Menu
	if err := database.DB.First(&menu, id).Error; err != nil {
		response.Error(c, "菜单不存在")
		return
	}

	if err := c.ShouldBindJSON(&menu); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Save(&menu).Error; err != nil {
		response.Error(c, "更新失败")
		return
	}

	response.SuccessWithMsg(c, "更新成功", menu)
}

func DeleteMenu(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var count int64
	database.DB.Model(&models.Menu{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		response.Error(c, "存在子菜单，无法删除")
		return
	}

	if err := database.DB.Delete(&models.Menu{}, id).Error; err != nil {
		response.Error(c, "删除失败")
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}
