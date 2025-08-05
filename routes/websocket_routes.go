package routes

import (
	"github.com/JZ23-2/splitbill-backend/handlers"
	"github.com/gin-gonic/gin"
)

func WebsocketRoutes(api *gin.RouterGroup) {
	websocket := api.Group("/ws")
	{
		websocket.GET("/:userId", handlers.WebSocketHandler)
	}
}
