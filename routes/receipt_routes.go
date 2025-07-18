package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func ReceiptRoutes(api *gin.RouterGroup) {
	receipt := api.Group("/receipt")
	{
		receipt.POST("/", controllers.ExtractReceipt)
	}
}