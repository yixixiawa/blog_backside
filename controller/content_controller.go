package controller

import (
	"sqlite_test/Model"
	"sqlite_test/constants"
	"sqlite_test/database"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateContent 创建内容
func CreateContent(c *gin.Context) {
	var content Model.Content
	if err := c.ShouldBindJSON(&content); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 从上下文中获取用户ID并正确转换类型
	userID, exists := c.Get("user_id")
	if !exists {
		constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}
	content.UserID = uint(userID.(int64))

	if err := database.DB.Create(&content).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, content)
}

// ListContents 获取内容列表
func ListContents(c *gin.Context) {
	var contents []Model.Content
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize
	query := database.DB.Offset(offset).Limit(pageSize)

	if err := query.Find(&contents).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, contents)
}

// GetContent 获取单个内容
func GetContent(c *gin.Context) {
	id := c.Param("id")
	var content Model.Content

	if err := database.DB.First(&content, id).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "内容不存在"})
		return
	}

	constants.SendResponse(c, constants.Success, content)
}

// UpdateContent 更新内容
func UpdateContent(c *gin.Context) {
	id := c.Param("id")
	var content Model.Content

	if err := database.DB.First(&content, id).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "内容不存在"})
		return
	}

	// 验证是否为作者
	userID, exists := c.Get("user_id")
	if !exists {
		constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}
	if content.UserID != uint(userID.(int64)) {
		constants.SendResponse(c, constants.UserForbidden, gin.H{"error": "无权修改他人内容"})
		return
	}

	if err := c.ShouldBindJSON(&content); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&content).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, content)
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
	userID, _ := c.Get("user_id")
	userID, exists := c.Get("user_id")
	if !exists {
		constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}
	if content.UserID != uint(userID.(int64)) {
		constants.SendResponse(c, constants.UserForbidden, gin.H{"error": "无权修改他人内容"})
		return
	}

	if err := database.DB.Delete(&content).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, nil)
}

// AddContentTags 添加内容标签
func AddContentTags(c *gin.Context) {
	// 获取并转换contentID
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
	var tagIDs []uint
	if err := c.ShouldBindJSON(&tagIDs); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 使用事务处理标签关联
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		for _, tagID := range tagIDs {
			// 验证标签是否存在
			var tag Model.Tag
			if err := tx.First(&tag, tagID).Error; err != nil {
				return err
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

	constants.SendResponse(c, constants.Success, nil)
}

// RemoveContentTag 移除内容标签
func RemoveContentTag(c *gin.Context) {
	contentID := c.Param("id")
	tagID := c.Param("tagId")

	if err := database.DB.Where("content_id = ? AND tag_id = ?", contentID, tagID).Delete(&Model.ContentTag{}).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, nil)
}
