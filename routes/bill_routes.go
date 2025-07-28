package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func BillRoutes(api *gin.RouterGroup) {
	bill := api.Group("/bills")
	{
		bill.POST("/", controllers.CreateBill)
		bill.POST("/bill-without-participant", controllers.CreateBillWithoutParticipant)
		bill.GET("/by-creator", controllers.GetBillByCreator)
	}
}
