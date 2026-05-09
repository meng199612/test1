package middleware

import (
	"strings"

	"admin-backend/internal/config"
	"admin-backend/pkg/jwt"
	"admin-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			response.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		jwtService := jwt.NewJWTService(config.AppConfig.JWT.Secret, config.AppConfig.JWT.Expire)
		claims, err := jwtService.ParseToken(parts[1])
		if err != nil {
			if err == jwt.ErrTokenExpired {
				response.Unauthorized(c, "登录已过期")
			} else {
				response.Unauthorized(c, "无效的token")
			}
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}
