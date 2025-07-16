package main

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/controllers"
	_ "github.com/JZ23-2/splitbill-backend/docs"
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
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.GET("api/v1/check", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "masuk"})
	})

	r.GET("api/v1/confirm-payment", controllers.ConfirmTransaction)

	r.Run(":8080")
}
