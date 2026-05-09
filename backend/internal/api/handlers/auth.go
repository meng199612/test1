package handlers

import (
	"admin-backend/internal/config"
	"admin-backend/internal/database"
	"admin-backend/internal/models"
	"admin-backend/pkg/jwt"
	"admin-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo interface{} `json:"user_info"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).Preload("Roles").First(&user).Error; err != nil {
		response.Error(c, "用户名或密码错误")
		return
	}

	if user.Status != 1 {
		response.Error(c, "账号已被禁用")
		return
	}

	if !user.CheckPassword(req.Password) {
		response.Error(c, "用户名或密码错误")
		return
	}

	jwtService := jwt.NewJWTService(config.AppConfig.JWT.Secret, config.AppConfig.JWT.Expire)
	token, err := jwtService.GenerateToken(user.ID, user.Username)
	if err != nil {
		response.Error(c, "生成token失败")
		return
	}

	response.Success(c, LoginResponse{
		Token: token,
		UserInfo: gin.H{
			"id":       user.ID,
			"username": user.Username,
			"nickname": user.Nickname,
			"avatar":   user.Avatar,
			"roles":    user.Roles,
		},
	})
}

func GetCurrentUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		response.Unauthorized(c, "请先登录")
		return
	}

	var user models.User
	if err := database.DB.Where("id = ?", userID).Preload("Roles").First(&user).Error; err != nil {
		response.Error(c, "获取用户信息失败")
		return
	}

	menuIDs, err := getUserMenuIDs(user.ID)
	if err != nil {
		menuIDs = []uint{}
	}

	permissionCodes, err := getUserPermissionCodes(user.ID)
	if err != nil {
		permissionCodes = []string{}
	}

	response.Success(c, gin.H{
		"id":         user.ID,
		"username":    user.Username,
		"nickname":    user.Nickname,
		"avatar":      user.Avatar,
		"roles":       user.Roles,
		"menu_ids":    menuIDs,
		"permissions": permissionCodes,
	})
}

func getUserMenuIDs(userID uint) ([]uint, error) {
	var roles []models.Role
	if err := database.DB.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Preload("Menus").
		Find(&roles).Error; err != nil {
		return nil, err
	}

	menuIDMap := make(map[uint]bool)
	for _, role := range roles {
		for _, menu := range role.Menus {
			menuIDMap[menu.ID] = true
		}
	}

	menuIDs := make([]uint, 0, len(menuIDMap))
	for id := range menuIDMap {
		menuIDs = append(menuIDs, id)
	}

	return menuIDs, nil
}

func getUserPermissionCodes(userID uint) ([]string, error) {
	var roles []models.Role
	if err := database.DB.Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Preload("Permissions").
		Find(&roles).Error; err != nil {
		return nil, err
	}

	codeMap := make(map[string]bool)
	for _, role := range roles {
		for _, perm := range role.Permissions {
			codeMap[perm.Code] = true
		}
	}

	var codes []string
	for code := range codeMap {
		codes = append(codes, code)
	}

	return codes, nil
}
