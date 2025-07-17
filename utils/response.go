package utils

import (
	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/gin-gonic/gin"
)

func SuccessResponse(c *gin.Context, status int, message string, data any) {
	c.JSON(status,
		dtos.SuccessResponse{
			Message: message,
			Data:    data,
		},
	)
}

func FailedResponse(c *gin.Context, status int, message string) {
	c.JSON(status, dtos.FailedResponse{
		Error: message,
	})
}
