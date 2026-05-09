package main

import (
	"fmt"
	"log"

	"admin-backend/internal/api/handlers"
	"admin-backend/internal/api/middleware"
	"admin-backend/internal/config"
	"admin-backend/internal/database"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	config.InitConfig()
	database.InitDB()

	if config.AppConfig.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.Use(middleware.LogMiddleware())

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", handlers.Login)
			auth.GET("/me", middleware.AuthMiddleware(), handlers.GetCurrentUser)
		}

		api.GET("/menus/tree", middleware.AuthMiddleware(), handlers.GetUserMenus)

		users := api.Group("/users")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("", handlers.GetUsers)
			users.GET("/:id", handlers.GetUser)
			users.POST("", handlers.CreateUser)
			users.PUT("/:id", handlers.UpdateUser)
			users.DELETE("/:id", handlers.DeleteUser)
			users.POST("/:id/reset-password", handlers.ResetPassword)
		}

		roles := api.Group("/roles")
		roles.Use(middleware.AuthMiddleware())
		{
			roles.GET("", handlers.GetRoles)
			roles.GET("/:id", handlers.GetRole)
			roles.POST("", handlers.CreateRole)
			roles.PUT("/:id", handlers.UpdateRole)
			roles.DELETE("/:id", handlers.DeleteRole)
			roles.POST("/:id/permissions", handlers.AssignPermissions)
			roles.POST("/:id/menus", handlers.AssignMenus)
		}

		permissions := api.Group("/permissions")
		permissions.Use(middleware.AuthMiddleware())
		{
			permissions.GET("/tree", handlers.GetPermissionsTree)
			permissions.POST("", handlers.CreatePermission)
			permissions.PUT("/:id", handlers.UpdatePermission)
			permissions.DELETE("/:id", handlers.DeletePermission)
		}

		menus := api.Group("/menus")
		menus.Use(middleware.AuthMiddleware())
		{
			menus.GET("/all", handlers.GetAllMenus)
			menus.GET("/:id", handlers.GetMenu)
			menus.POST("", handlers.CreateMenu)
			menus.PUT("/:id", handlers.UpdateMenu)
			menus.DELETE("/:id", handlers.DeleteMenu)
		}

		logs := api.Group("/logs")
		logs.Use(middleware.AuthMiddleware())
		{
			logs.GET("", handlers.GetOperationLogs)
			logs.GET("/:id", handlers.GetOperationLog)
			logs.DELETE("/:id", handlers.DeleteOperationLog)
			logs.POST("/clear", handlers.ClearOperationLogs)
		}

		upload := api.Group("/upload")
		upload.Use(middleware.AuthMiddleware())
		{
			upload.POST("", handlers.UploadFile)
			upload.POST("/multi", handlers.UploadFiles)
		}

		api.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "pong"})
		})
	}

	addr := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	log.Printf("Server starting on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
