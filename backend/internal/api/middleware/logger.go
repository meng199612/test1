package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"time"

	"admin-backend/internal/database"
	"admin-backend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		w := &responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBufferString(""),
		}
		c.Writer = w

		c.Next()

		duration := time.Since(startTime)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path
		clientIP := c.ClientIP()

		responseBody := w.body.String()

		var respMap map[string]interface{}
		logStatus := int8(1)
		if err := json.Unmarshal([]byte(responseBody), &respMap); err == nil {
			if code, ok := respMap["code"].(float64); ok && code != 200 {
				logStatus = 0
			}
		}

		logrus.WithFields(logrus.Fields{
			"status":   status,
			"duration": duration,
			"method":   method,
			"path":     path,
			"ip":       clientIP,
		}).Info("HTTP Request")

		if database.DB != nil && path != "/api/auth/login" {
			module := getModule(path)
			operation := getOperation(method, path)

			userID, _ := c.Get("user_id")
			username, _ := c.Get("username")

			var uid *uint
			if id, ok := userID.(uint); ok {
				uid = &id
			}

			uname := ""
			if u, ok := username.(string); ok {
				uname = u
			}

			logEntry := models.OperationLog{
				UserID:    uid,
				Username:  uname,
				Module:    module,
				Operation: operation,
				Method:    method,
				Path:      path,
				IP:        clientIP,
				Params:    string(requestBody),
				Result:    responseBody,
				Status:    logStatus,
				CreatedAt: time.Now(),
			}

			go func() {
				database.DB.Create(&logEntry)
			}()
		}
	}
}

func getModule(path string) string {
	if len(path) < 5 {
		return "other"
	}
	p := path[5:]
	for i, ch := range p {
		if ch == '/' {
			return p[:i]
		}
	}
	return p
}

func getOperation(method, path string) string {
	switch method {
	case "GET":
		return "查询"
	case "POST":
		return "新增"
	case "PUT":
		return "更新"
	case "DELETE":
		return "删除"
	default:
		return "其他"
	}
}
