package controllers

import "github.com/gin-gonic/gin"

// CheckBackendHealth godoc
// @Summary Check backend health
// @Description Check backend health
// @Tags Payment
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /check [get]
func CheckHealth(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Sehat Bro!",
	})
}
