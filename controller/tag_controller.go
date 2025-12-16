package controller

import (
	"sqlite_test/Model"
	"sqlite_test/constants"
	"sqlite_test/database"

	"github.com/gin-gonic/gin"
)

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var tag Model.Tag
	if err := c.ShouldBindJSON(&tag); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Create(&tag).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, tag)
}

func GetTag(c *gin.Context) {
	id := c.Param("id")
	var tag Model.Tag
	if err := database.DB.First(&tag, id).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "标签不存在"})
		return
	}
	constants.SendResponse(c, constants.Success, tag)
}

// ListTags 获取标签列表
func ListTags(c *gin.Context) {
	// 定义一个包含计数的临时结构体
	type TagWithCount struct {
		Model.Tag
		ArticleCount int64 `json:"article_count"`
	}
	var tags []TagWithCount

	// 联表查询并统计数量
	if err := database.DB.Table("tags").
		Select("tags.*, count(content_tags.content_id) as article_count").
		Joins("LEFT JOIN content_tags ON content_tags.tag_id = tags.tag_id").
		Group("tags.tag_id").
		Scan(&tags).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, tags)
}

// UpdateTag 更新标签
func UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var tag Model.Tag

	if err := database.DB.First(&tag, id).Error; err != nil {
		constants.SendResponse(c, constants.UserNotFound, gin.H{"error": "标签不存在"})
		return
	}

	if err := c.ShouldBindJSON(&tag); err != nil {
		constants.SendResponse(c, constants.UserBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := database.DB.Save(&tag).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, tag)
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&Model.Tag{}, id).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, nil)
}
