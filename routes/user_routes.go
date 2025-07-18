package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func UserRoutes(api *gin.RouterGroup) {
	user := api.Group("/users")
	{
		user.POST("/", controllers.RegisterUser)
	}
}
