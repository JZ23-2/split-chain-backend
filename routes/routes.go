package routes

import (
	"fmt"

	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetUpRoutes() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "PATCH", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{
		r.GET("api/v1/check", controllers.CheckHealth)

		UserRoutes(api)
		fmt.Println("logging here")
		// ParticipantRoutes(api)
		BillRoutes(api)
		ReceiptRoute(api)
		ExchangeRateRoute(api)
		FriendRoutes(api)
		PaymentRoutes(api)
		WebsocketRoutes(api)
		InboxRoutes(api)
	}

	r.Run(":8080")
}
