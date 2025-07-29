package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func BillRoutes(api *gin.RouterGroup) {
	bill := api.Group("/bills")
	{
		bill.POST("/assign-participants", controllers.AssignParticipantsController)
		bill.POST("/bill-without-participant", controllers.CreateBillWithoutParticipantController)
		bill.GET("/by-creator", controllers.GetBillByCreatorController)
	}
}
