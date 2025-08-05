package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func PaymentRoutes(api *gin.RouterGroup) {
	payment := api.Group("/payments")
	{
		payment.GET("/confirm-payment/:transactionId", controllers.ConfirmTransaction)
	}
}
