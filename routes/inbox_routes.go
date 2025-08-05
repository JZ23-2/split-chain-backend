package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func InboxRoutes(api *gin.RouterGroup) {
	inbox := api.Group("/inboxes")
	{
		inbox.GET("", controllers.GetUnreadInboxesByUserIDHandler)
		inbox.PATCH("/:id/read", controllers.MarkInboxAsReadHandler)
	}
}
