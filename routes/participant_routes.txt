package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func ParticipantRoutes(api *gin.RouterGroup) {
	participant := api.Group("/participants")
	{
		participant.GET("/:participant_id", controllers.GetParticipantBills)
		participant.GET("/:participant_id/:bill_id", controllers.GetParticipantDetail)
	}
}
