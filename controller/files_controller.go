package controller

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"sqlite_test/Model"
	"sqlite_test/constants"
	"sqlite_test/database"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetImageStoragePath 获取图片存储目录路径
func GetImageStoragePath() string {
	dir, _ := os.Getwd()
	imgDir := filepath.Join(dir, "img")
	_ = os.MkdirAll(imgDir, 0755)
	return imgDir
}

// randSuffix 生成随机后缀
func randSuffix(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)[:n]
}

// UploadFile 上传图片文件（仅支持 png/jpg/jpeg，最大 10MB）
func UploadFile(c *gin.Context) {
	const maxSize = 10 << 20 // 10MB
	allowedExt := map[string]bool{".png": true, ".jpg": true, ".jpeg": true}

	// 获取上传文件
	fh, err := c.FormFile("file")
	if err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{
			"error": "file is required",
			"hint":  "use form-data with field name 'file'",
		})
		return
	}

	// 校验文件大小
	if fh.Size > maxSize {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{
			"error":     "file too large",
			"max_bytes": maxSize,
		})
		return
	}

	// 校验文件扩展名
	ext := strings.ToLower(filepath.Ext(fh.Filename))
	if !allowedExt[ext] {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{
			"error":   "invalid file type",
			"allowed": []string{"png", "jpg", "jpeg"},
		})
		return
	}

	// 生成唯一文件名
	storageName := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), randSuffix(6), ext)
	dst := filepath.Join(GetImageStoragePath(), storageName)

	// 保存文件到磁盘
	if err := c.SaveUploadedFile(fh, dst); err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{
			"error":  "save file failed",
			"detail": err.Error(),
		})
		return
	}

	// 创建文件记录
	fr := Model.FileRecord{
		OriginalName: fh.Filename,
		StorageName:  storageName,
		FilePath:     dst,
		FileURL:      "/img/" + storageName,
		FileSize:     fh.Size,
		UploadTime:   time.Now(),
		Status:       "active",
	}

	// 使用事务：先保存文件记录，再创建关联（如果提供了 content_id）
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 保存文件记录
		if err := tx.Create(&fr).Error; err != nil {
			return err
		}

		// 如果提供了 content_id，创建关联关系
		if cid := c.PostForm("content_id"); cid != "" {
			contentID, err := strconv.ParseUint(cid, 10, 32)
			if err != nil {
				return fmt.Errorf("invalid content_id")
			}

			// 验证内容是否存在
			var content Model.Content
			if err := tx.First(&content, contentID).Error; err != nil {
				return fmt.Errorf("content not found")
			}

			// 获取当前最大 order
			var maxOrder int
			tx.Model(&Model.ContentFile{}).
				Where("content_id = ?", contentID).
				Select("COALESCE(MAX(`order`), -1)").
				Scan(&maxOrder)

			// 创建文章-文件关联
			contentFile := Model.ContentFile{
				ContentID: uint(contentID),
				FileID:    fr.FileID,
				Order:     maxOrder + 1,
				Usage:     c.DefaultPostForm("usage", "content"), // 可选参数
			}
			if err := tx.Create(&contentFile).Error; err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		_ = os.Remove(dst) // 删除已保存的文件
		constants.SendResponse(c, constants.UserSystemError, gin.H{
			"error":  "store record failed",
			"detail": err.Error(),
		})
		return
	}
	constants.SendResponse(c, constants.Success, fr)
}

// ListImages 获取所有图片列表（从数据库查询）
func ListImages(c *gin.Context) {
	var records []Model.FileRecord

	// 从数据库查询所有图片记录
	if err := database.DB.Where("status = ?", "active").Find(&records).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{
			"error":  "query images failed",
			"detail": err.Error(),
		})
		return
	}

	constants.SendResponse(c, constants.Success, gin.H{
		"total":  len(records),
		"images": records,
	})
}

// GetContentImages 根据文章ID获取关联的图片列表
func GetContentImages(c *gin.Context) {
	contentID := c.Param("id")

	var contentFiles []Model.ContentFile
	if err := database.DB.Preload("FileRecord").
		Where("content_id = ?", contentID).
		Order("`order` ASC, created_at ASC").
		Find(&contentFiles).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{
			"error":  "query content images failed",
			"detail": err.Error(),
		})
		return
	}

	// 提取文件记录列表
	var images []Model.FileRecord
	for _, cf := range contentFiles {
		images = append(images, cf.FileRecord)
	}

	constants.SendResponse(c, constants.Success, gin.H{
		"content_id": contentID,
		"total":      len(images),
		"images":     images,
	})
}
