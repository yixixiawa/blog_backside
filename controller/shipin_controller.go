package controller

import (
	"sqlite_test/Model"
	"sqlite_test/constants"
	"sqlite_test/database"
	"strconv"

	"github.com/gin-gonic/gin"
)

func search_goods(c *gin.Context) {
	var goods Model.ItemQueryRequest
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	offset := (page - 1) * pageSize
	query := database.DB.Offset(offset).Limit(pageSize)

	if err := query.Find(&goods).Error; err != nil {
		constants.SendResponse(c, constants.UserSystemError, gin.H{"error": err.Error()})
		return
	}

	constants.SendResponse(c, constants.Success, goods)
}
