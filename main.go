package main

import (
	"github.com/JZ23-2/splitbill-backend/config"
	"github.com/JZ23-2/splitbill-backend/controllers"
	"github.com/JZ23-2/splitbill-backend/database"
	_ "github.com/JZ23-2/splitbill-backend/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @Split Chain
// @version 1.0
// @Split Chain Backend
// @contact.Jackson API Support
// @contact.email Jacksontpa7@gmail.com
// @license.name MIT
// @BasePath /api/v1
func main() {
	config.Loadenv()
	database.ConnectDB()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api/v1")
	{

		r.GET("api/v1/check", controllers.CheckHealth)
		r.GET("api/v1/confirm-payment", controllers.ConfirmTransaction)

		user := api.Group("/users")
		{
			user.POST("/register", controllers.RegisterUser)
		}

		participant := api.Group("/participants")
		{
			participant.GET("/get-all-participant-detail/:participantId", controllers.GetParticipantBills)

			participant.POST("/get-participant-detail", controllers.GetParticipantDetail)
		}

		bill := api.Group("/bills")
		{
			bill.POST("/create-bill", controllers.CreateBill)
		}

	}

	r.Run(":8080")

}
