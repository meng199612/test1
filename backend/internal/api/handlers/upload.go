package handlers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"admin-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		response.BadRequest(c, "文件大小不能超过10MB")
		return
	}

	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
		".pdf": true, ".doc": true, ".docx": true, ".xls": true, ".xlsx": true,
		".txt": true, ".zip": true, ".rar": true,
	}
	if !allowedExts[ext] {
		response.BadRequest(c, "不支持的文件类型")
		return
	}

	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.Error(c, "创建上传目录失败")
		return
	}

	timestamp := time.Now().Unix()
	newFilename := fmt.Sprintf("%d%s", timestamp, ext)
	filePath := filepath.Join(uploadDir, newFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		response.Error(c, "创建文件失败")
		return
	}
	defer dst.Close()

	src, err := file.Open()
	if err != nil {
		response.Error(c, "读取文件失败")
		return
	}
	defer src.Close()

	if _, err = io.Copy(dst, src); err != nil {
		response.Error(c, "保存文件失败")
		return
	}

	host := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	fileURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, host, newFilename)

	response.Success(c, gin.H{
		"url":  fileURL,
		"name": file.Filename,
		"size": file.Size,
	})
}

func GetUploadedFile(c *gin.Context) {
	filename := c.Param("filename")
	filePath := filepath.Join("./uploads", filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(404, gin.H{"error": "文件不存在"})
		return
	}

	c.File(filePath)
}

func UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.BadRequest(c, "请选择要上传的文件")
		return
	}

	uploadDir := "./uploads"
	os.MkdirAll(uploadDir, 0755)

	var results []gin.H
	host := c.Request.Host
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}

	for _, file := range files {
		if file.Size > int64(10*1024*1024) {
			continue
		}

		ext := filepath.Ext(file.Filename)
		timestamp := time.Now().UnixNano()
		newFilename := fmt.Sprintf("%d%s", timestamp, ext)
		filePath := filepath.Join(uploadDir, newFilename)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			continue
		}

		fileURL := fmt.Sprintf("%s://%s/uploads/%s", scheme, host, newFilename)
		results = append(results, gin.H{
			"url":  fileURL,
			"name": file.Filename,
			"size": file.Size,
		})
	}

	response.Success(c, gin.H{"files": results})
}
