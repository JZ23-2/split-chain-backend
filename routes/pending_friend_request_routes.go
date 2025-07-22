package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func PendingFriendRoutes(api *gin.RouterGroup) {
	friend_request := api.Group("/pending-friend-request")
	{
		friend_request.POST("/send-request", controllers.SendFriendRequest)
	}
}
