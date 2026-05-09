package handlers

import (
	"strconv"
	"time"

	"admin-backend/internal/database"
	"admin-backend/internal/models"
	"admin-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func GetOperationLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	username := c.Query("username")
	module := c.Query("module")
	startTime := c.Query("start_time")
	endTime := c.Query("end_time")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	query := database.DB.Model(&models.OperationLog{})
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}
	if startTime != "" {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != "" {
		query = query.Where("created_at <= ?", endTime+" 23:59:59")
	}

	var total int64
	query.Count(&total)

	var logs []models.OperationLog
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&logs)

	response.Page(c, logs, total, page, pageSize)
}

func GetOperationLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	var log models.OperationLog
	if err := database.DB.First(&log, id).Error; err != nil {
		response.Error(c, "日志不存在")
		return
	}

	response.Success(c, log)
}

func DeleteOperationLog(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	if err := database.DB.Delete(&models.OperationLog{}, id).Error; err != nil {
		response.Error(c, "删除失败")
		return
	}

	response.SuccessWithMsg(c, "删除成功", nil)
}

func ClearOperationLogs(c *gin.Context) {
	var req struct {
		Days int `json:"days"`
	}
	c.ShouldBindJSON(&req)

	if req.Days <= 0 {
		req.Days = 30
	}

	cutoff := time.Now().AddDate(0, 0, -req.Days)
	result := database.DB.Where("created_at < ?", cutoff).Delete(&models.OperationLog{})
	if result.Error != nil {
		response.Error(c, "清理失败")
		return
	}

	response.SuccessWithMsg(c, "清理成功", gin.H{"deleted": result.RowsAffected})
}
