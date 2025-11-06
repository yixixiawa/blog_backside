package constants

import "github.com/gin-gonic/gin"

func SendResponse(c *gin.Context, status StatusCode, data interface{}) {
	c.JSON(status.GetCode(), BuildResponseWithStatus(status, data))
}

// 每个控制层的独立包装方法，方便控制层直接调用（保持类型安全与明确性）
func SendUserResponse(c *gin.Context, status UserStatusCode, data interface{}) {
	c.JSON(status.GetCode(), BuildResponseWithStatus(status, data))
}

func SendContentResponse(c *gin.Context, status ContentStatusCode, data interface{}) {
	c.JSON(status.GetCode(), BuildResponseWithStatus(status, data))
}

func SendCommentResponse(c *gin.Context, status CommentStatusCode, data interface{}) {
	c.JSON(status.GetCode(), BuildResponseWithStatus(status, data))
}

func SendTagResponse(c *gin.Context, status TagStatusCode, data interface{}) {
	c.JSON(status.GetCode(), BuildResponseWithStatus(status, data))
}

func SendEmailResponse(c *gin.Context, status EmailStatusCode, data interface{}) {
	c.JSON(status.GetCode(), BuildResponseWithStatus(status, data))
}
