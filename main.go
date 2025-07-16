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

	r.GET("api/v1/check", controllers.CheckHealth)

	r.GET("api/v1/confirm-payment", controllers.ConfirmTransaction)

	r.POST("api/v1/register", controllers.RegisterUser)

	r.Run(":8080")
}
