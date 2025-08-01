package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func BillRoutes(api *gin.RouterGroup) {
	bill := api.Group("/bills")
	{
		bill.POST("/assign-participants", controllers.AssignParticipantController)
		bill.POST("/bill-without-participant", controllers.CreateBillWithoutParticipantController)
		bill.GET("/by-creator", controllers.GetBillByCreatorController)
		bill.GET("/by-participant/:participantId", controllers.GetBillsByParticipantController)
		bill.GET("/by-billId/:billId", controllers.GetBillByBillIDHandler)
		bill.DELETE("/delete-bill/:billId", controllers.DeleteBillByIDController)
		bill.PATCH("/update-bill", controllers.UpdateBillController)
	}
}
