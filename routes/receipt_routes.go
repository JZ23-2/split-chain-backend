package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func ReceiptRoute(api *gin.RouterGroup) {
	api.POST("/receipt", controllers.ExtractReceipt)
}
