package routes

import (
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-gonic/gin"
)

func ExchangeRateRoute(api *gin.RouterGroup) {
	api.GET("/get-rate", controllers.ConvertPricesToHBAR)
}
