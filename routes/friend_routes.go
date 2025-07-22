package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func FriendRoutes(api *gin.RouterGroup) {
	friend := api.Group("/friend")
	{
		friend.POST("/accept-friend", controllers.AcceptFriendRequest)
		friend.POST("/decline-friend", controllers.DeclineFriendRequest)
	}
}
