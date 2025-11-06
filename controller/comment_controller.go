package controller

import (
	"sqlite_test/Model"
	"sqlite_test/constants"
	"sqlite_test/database"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateComment 创建评论
func CreateComment(c *gin.Context) {
	var comment Model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 设置评论者ID
	userID, exists := c.Get("user_id")
	if !exists {
		constants.SendResponse(c, constants.UserUnauthorized, gin.H{"error": "未获取到用户信息"})
		return
	}
	comment.UserID = uint(userID.(int64))

	// 验证内容是否存在
	var content Model.Content
	if err := database.DB.First(&content, comment.ContentID).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "评论的内容不存在"})
		return
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, comment)
}

// ListContentComments 获取内容的评论列表
func ListContentComments(c *gin.Context) {
	contentIDStr := c.Param("contentId")
	contentID, err := strconv.ParseUint(contentIDStr, 10, 32)
	if err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": "无效的内容ID"})
		return
	}

	var comments []Model.Comment
	// 预加载用户信息，但不返回用户敏感信息
	if err := database.DB.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("user_id", "username")
	}).Where("content_id = ?", uint(contentID)).Find(&comments).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, comments)
}

// UpdateComment 更新评论
func UpdateComment(c *gin.Context) {
	id := c.Param("id")
	var existingComment Model.Comment

	if err := database.DB.First(&existingComment, id).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 验证是否为评论作者
	userID, _ := c.Get("user_id")
	if existingComment.UserID != uint(userID.(int64)) {
		constants.SendResponse(c, constants.UserForbidden, gin.H{"error": "无权修改他人评论"})
		return
	}

	// 只更新评论内容
	var updateData struct {
		CommentText string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&updateData); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingComment.CommentText = updateData.CommentText
	if err := database.DB.Save(&existingComment).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, existingComment)
}

// DeleteComment 删除评论
func DeleteComment(c *gin.Context) {
	id := c.Param("id")
	var comment Model.Comment

	if err := database.DB.First(&comment, id).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 验证是否为评论作者
	userID, _ := c.Get("user_id")
	if comment.UserID != uint(userID.(int64)) {
		constants.SendResponse(c, constants.UserForbidden, gin.H{"error": "无权删除他人评论"})
		return
	}

	if err := database.DB.Delete(&comment).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, nil)
}
