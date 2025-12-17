package controller

import (
	"blog/Model"
	"blog/constants"
	"blog/database"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// getUserID 统一获取用户ID的辅助函数
func getUserID(c *gin.Context) (uint, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return 0, fmt.Errorf("未获取到用户信息")
	}
	return uint(userID.(int64)), nil
}

// CreateContent 创建内容
func CreateContent(c *gin.Context) {
	var content Model.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取用户ID
	userID, err := getUserID(c)
	if err != nil {
		constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": err.Error()})
		return
	}
	content.UserID = userID

	if err := database.DB.Create(&content).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, content)
}

// ListContents 获取内容列表（增加总数和预加载）
func ListContents(c *gin.Context) {
	var contents []Model.Content
	var total int64

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// 参数验证
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	// 获取总数
	database.DB.Model(&Model.Content{}).Count(&total)

	// 查询数据，预加载关联
	if err := database.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "username", "avatar")
		}).
		Preload("Tags").
		Offset(offset).
		Limit(pageSize).
		Order("created_at DESC").
		Find(&contents).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, gin.H{
		"list":      contents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetContent 获取单个内容（增加预加载）
func GetContent(c *gin.Context) {
	id := c.Param("id")
	var content Model.Content

	if err := database.DB.
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "username", "avatar")
		}).
		Preload("Tags").
		Preload("Comments.User", func(db *gorm.DB) *gorm.DB {
			return db.Select("user_id", "username", "avatar")
		}).
		Preload("ContentFiles", func(db *gorm.DB) *gorm.DB {
			return db.Order("`order` ASC")
		}).
		Preload("ContentFiles.FileRecord").
		First(&content, id).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "内容不存在"})
		return
	}

	// 增加浏览量
	database.DB.Model(&content).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	constants.SendResponse(c, constants.Success, content)
}

// UpdateContent 更新内容
func UpdateContent(c *gin.Context) {
	id := c.Param("id")
	var existingContent Model.Content

	if err := database.DB.First(&existingContent, id).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "内容不存在"})
		return
	}

	// 验证是否为作者
	userID, err := getUserID(c)
	if err != nil {
		constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if existingContent.UserID != userID {
		constants.SendResponse(c, constants.UserForbidden, gin.H{"error": "无权修改他人内容"})
		return
	}

	// 只更新允许修改的字段
	var updateData struct {
		Title             string `json:"title"`
		Content           string `json:"content"`
		BriefIntroduction string `json:"brief_introduction"`
		CoverImage        string `json:"cover_image"`
		Status            string `json:"status"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 更新字段
	updates := map[string]interface{}{
		"title":              updateData.Title,
		"content":            updateData.Content,
		"brief_introduction": updateData.BriefIntroduction,
		"cover_image":        updateData.CoverImage,
		"status":             updateData.Status,
	}

	if err := database.DB.Model(&existingContent).Updates(updates).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, existingContent)
}

// DeleteContent 删除内容
func DeleteContent(c *gin.Context) {
	id := c.Param("id")
	var content Model.Content

	if err := database.DB.First(&content, id).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "内容不存在"})
		return
	}

	// 验证是否为作者
	userID, err := getUserID(c)
	if err != nil {
		constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if content.UserID != userID {
		constants.SendResponse(c, constants.UserForbidden, gin.H{"error": "无权删除他人内容"})
		return
	}

	// 使用事务删除内容及其关联
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// 删除标签关联
		if err := tx.Where("content_id = ?", id).Delete(&Model.ContentTag{}).Error; err != nil {
			return err
		}
		// 删除文件关联
		if err := tx.Where("content_id = ?", id).Delete(&Model.ContentFile{}).Error; err != nil {
			return err
		}
		// 删除评论
		if err := tx.Where("content_id = ?", id).Delete(&Model.Comment{}).Error; err != nil {
			return err
		}
		// 删除内容
		if err := tx.Delete(&content).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, nil)
}

// AddContentTags 添加内容标签
func AddContentTags(c *gin.Context) {
	contentIDStr := c.Param("id")
	contentID, err := strconv.ParseUint(contentIDStr, 10, 32)
	if err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": "无效的内容ID"})
		return
	}

	// 验证内容是否存在
	var content Model.Content
	if err := database.DB.First(&content, contentID).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "内容不存在"})
		return
	}

	// 获取标签ID列表
	var req struct {
		TagIDs []uint `json:"tag_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用事务处理标签关联
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		for _, tagID := range req.TagIDs {
			// 检查标签是否已关联
			var count int64
			tx.Model(&Model.ContentTag{}).
				Where("content_id = ? AND tag_id = ?", contentID, tagID).
				Count(&count)

			if count > 0 {
				continue // 已存在，跳过
			}

			// 验证标签是否存在
			var tag Model.Tag
			if err := tx.First(&tag, tagID).Error; err != nil {
				return fmt.Errorf("标签ID %d 不存在", tagID)
			}

			contentTag := Model.ContentTag{
				ContentID: uint(contentID),
				TagID:     tagID,
			}
			if err := tx.Create(&contentTag).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, gin.H{"message": "标签添加成功"})
}

// RemoveContentTag 移除内容标签
func RemoveContentTag(c *gin.Context) {
	contentID := c.Param("id")
	tagID := c.Param("tagId")

	result := database.DB.Where("content_id = ? AND tag_id = ?", contentID, tagID).Delete(&Model.ContentTag{})
	if result.Error != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "未找到该标签关联"})
		return
	}

	constants.SendResponse(c, constants.Success, gin.H{"message": "标签移除成功"})
}

// AddContentFiles 为内容添加文件关联
func AddContentFiles(c *gin.Context) {
	contentID := c.Param("id")

	var req struct {
		FileIDs []uint `json:"file_ids" binding:"required"`
		Usage   string `json:"usage"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 验证内容是否存在
	var content Model.Content
	if err := database.DB.First(&content, contentID).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "内容不存在"})
		return
	}

	// 使用事务批量插入
	err := database.DB.Transaction(func(tx *gorm.DB) error {
		// 获取当前最大 order 值
		var maxOrder int
		tx.Model(&Model.ContentFile{}).
			Where("content_id = ?", contentID).
			Select("COALESCE(MAX(`order`), -1)").
			Scan(&maxOrder)

		for i, fileID := range req.FileIDs {
			// 验证文件是否存在
			var file Model.FileRecord
			if err := tx.First(&file, fileID).Error; err != nil {
				return fmt.Errorf("文件ID %d 不存在", fileID)
			}

			contentFile := Model.ContentFile{
				ContentID: content.ID,
				FileID:    fileID,
				Usage:     req.Usage,
				Order:     maxOrder + i + 1,
			}
			if err := tx.Create(&contentFile).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, gin.H{"message": "文件关联成功"})
}

// RemoveContentFile 移除内容的文件关联
func RemoveContentFile(c *gin.Context) {
	contentID := c.Param("id")
	fileID := c.Param("fileId")

	result := database.DB.Where("content_id = ? AND file_id = ?", contentID, fileID).Delete(&Model.ContentFile{})
	if result.Error != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": result.Error.Error()})
		return
	}
	if result.RowsAffected == 0 {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "未找到文件关联"})
		return
	}

	constants.SendResponse(c, constants.Success, gin.H{"message": "文件关联移除成功"})
}

// GetContentFiles 获取内容关联的所有文件
func GetContentFiles(c *gin.Context) {
	contentID := c.Param("id")

	var contentFiles []Model.ContentFile
	if err := database.DB.Preload("FileRecord").
		Where("content_id = ?", contentID).
		Order("`order` ASC").
		Find(&contentFiles).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, contentFiles)
}
