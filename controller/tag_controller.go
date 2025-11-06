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
	var tags []Model.Tag

	if err := database.DB.Find(&tags).Error; err != nil {
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
