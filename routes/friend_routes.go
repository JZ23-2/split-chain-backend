package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func FriendRoutes(api *gin.RouterGroup) {
	friend := api.Group("/friends")
	{
		friend.POST("/accept", controllers.AcceptFriendRequest)
		friend.POST("/decline", controllers.DeclineFriendRequest)
		friend.POST("/add", controllers.AddFriend)
		friend.GET("/:user_wallet_address", controllers.GetFriend)
	}
}
