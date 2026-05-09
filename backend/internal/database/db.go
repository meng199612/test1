package database

import (
	"log"
	"time"

	"admin-backend/internal/config"
	"admin-backend/internal/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := config.AppConfig.Database.DSN()

	var err error
	for i := 0; i < 30; i++ {
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
		if err == nil {
			sqlDB, err := DB.DB()
			if err == nil {
				if err = sqlDB.Ping(); err == nil {
					log.Printf("Database connected on attempt %d", i+1)
					break
				}
			}
		}
		log.Printf("Waiting for database (attempt %d/30): %v", i+1, err)
		time.Sleep(2 * time.Second)
	}

	if DB == nil {
		log.Fatalf("Failed to connect to database after 30 attempts")
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.Menu{},
		&models.OperationLog{},
	)
	if err != nil {
		log.Fatalf("Auto migrate failed: %v", err)
	}

	initDefaultData()

	log.Println("Database initialized successfully")
}

func initDefaultData() {
	var count int64
	DB.Model(&models.Role{}).Count(&count)
	if count > 0 {
		return
	}

	roles := []models.Role{
		{Name: "超级管理员", Code: "super_admin", Description: "系统最高权限管理员"},
		{Name: "普通管理员", Code: "admin", Description: "普通管理员"},
		{Name: "普通用户", Code: "user", Description: "普通用户"},
	}
	DB.Create(&roles)

	menus := []models.Menu{
		{Name: "系统管理", Path: "/system", Icon: "setting", ParentID: 0, Component: "", Sort: 1},
		{Name: "用户管理", Path: "/system/users", Icon: "user", ParentID: 1, Component: "system/users/index", Sort: 1},
		{Name: "角色管理", Path: "/system/roles", Icon: "user-role", ParentID: 1, Component: "system/roles/index", Sort: 2},
		{Name: "菜单管理", Path: "/system/menus", Icon: "menu", ParentID: 1, Component: "system/menus/index", Sort: 3},
		{Name: "日志管理", Path: "/system/logs", Icon: "document", ParentID: 1, Component: "system/logs/index", Sort: 4},
	}
	DB.Create(&menus)

	permissions := []models.Permission{
		{Name: "用户管理", Code: "user:manage", Type: 1, ParentID: 0, Path: "/api/users", Method: "GET", Sort: 1},
		{Name: "用户列表", Code: "user:list", Type: 2, ParentID: 1, Path: "/api/users", Method: "GET", Sort: 1},
		{Name: "新增用户", Code: "user:add", Type: 2, ParentID: 1, Path: "/api/users", Method: "POST", Sort: 2},
		{Name: "编辑用户", Code: "user:edit", Type: 2, ParentID: 1, Path: "/api/users/:id", Method: "PUT", Sort: 3},
		{Name: "删除用户", Code: "user:delete", Type: 2, ParentID: 1, Path: "/api/users/:id", Method: "DELETE", Sort: 4},
		{Name: "角色管理", Code: "role:manage", Type: 1, ParentID: 0, Path: "/api/roles", Method: "GET", Sort: 2},
		{Name: "角色列表", Code: "role:list", Type: 2, ParentID: 6, Path: "/api/roles", Method: "GET", Sort: 1},
		{Name: "新增角色", Code: "role:add", Type: 2, ParentID: 6, Path: "/api/roles", Method: "POST", Sort: 2},
		{Name: "编辑角色", Code: "role:edit", Type: 2, ParentID: 6, Path: "/api/roles/:id", Method: "PUT", Sort: 3},
		{Name: "删除角色", Code: "role:delete", Type: 2, ParentID: 6, Path: "/api/roles/:id", Method: "DELETE", Sort: 4},
		{Name: "分配权限", Code: "role:permission", Type: 2, ParentID: 6, Path: "/api/roles/:id/permissions", Method: "POST", Sort: 5},
	}
	DB.Create(&permissions)

	admin := models.User{
		Username: "admin",
		Nickname: "超级管理员",
		Status:   1,
	}
	admin.SetPassword("admin123")
	DB.Create(&admin)
	DB.Model(&admin).Association("Roles").Append(&roles[0])

	log.Println("Default data initialized")
}
