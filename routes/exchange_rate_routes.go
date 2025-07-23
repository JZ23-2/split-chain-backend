package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func ExchangeRateRoute(api *gin.RouterGroup) {
	api.POST("/rate", controllers.ConvertPricesToHBAR)
}
