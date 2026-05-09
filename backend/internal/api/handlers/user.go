package handlers

import (
	"strconv"

	"admin-backend/internal/database"
	"admin-backend/internal/models"
	"admin-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int8   `json:"status"`
	RoleIDs  []uint `json:"role_ids"`
}

type UpdateUserRequest struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Status   int8   `json:"status"`
	RoleIDs  []uint `json:"role_ids"`
}

func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	username := c.Query("username")
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.DB.Model(&models.User{}).Preload("Roles")
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var users []models.User
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&users)

	response.Page(c, users, total, page, pageSize)
}

func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.Preload("Roles").First(&user, id).Error; err != nil {
		response.Error(c, "用户不存在")
		return
	}

	response.Success(c, user)
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	var existing models.User
	if database.DB.Where("username = ?", req.Username).First(&existing).Error == nil {
		response.Error(c, "用户名已存在")
		return
	}

	user := models.User{
		Username: req.Username,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   req.Status,
	}

	if err := user.SetPassword(req.Password); err != nil {
		response.Error(c, "密码加密失败")
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		response.Error(c, "创建用户失败")
		return
	}

	if len(req.RoleIDs) > 0 {
		var roles []models.Role
		database.DB.Where("id IN ?", req.RoleIDs).Find(&roles)
		database.DB.Model(&user).Association("Roles").Replace(&roles)
	}

	response.SuccessWithMsg(c, "创建成功", user)
}

func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		response.Error(c, "用户不存在")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Nickname != "" {
		updates["nickname"] = req.Nickname
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Status != 0 || user.Status != req.Status {
		updates["status"] = req.Status
	}

	if len(updates) > 0 {
		if err := database.DB.Model(&user).Updates(updates).Error; err != nil {
			response.Error(c, "更新失败")
			return
		}
	}

	if req.RoleIDs != nil {
		var roles []models.Role
		database.DB.Where("id IN ?", req.RoleIDs).Find(&roles)
		database.DB.Model(&user).Association("Roles").Replace(&roles)
	}

	response.SuccessWithMsg(c, "更新成功", user)
}

func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		response.Error(c, "用户不存在")
		return
	}

	if user.Username == "admin" {
		response.Error(c, "超级管理员不能删除")
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		response.Error(c, "删除失败")
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

func ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Password string `json:"password" binding:"required,min=6"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		response.Error(c, "用户不存在")
		return
	}

	if err := user.SetPassword(req.Password); err != nil {
		response.Error(c, "密码加密失败")
		return
	}

	if err := database.DB.Save(&user).Error; err != nil {
		response.Error(c, "重置密码失败")
		return
	}

	response.SuccessWithMsg(c, "重置密码成功", nil)
}
