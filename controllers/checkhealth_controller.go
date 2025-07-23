package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

// CheckBackendHealth godoc
//	@Summary		Check backend health
//	@Description	Check backend health
//	@Tags			Payment
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/check [get]
func CheckHealth(c *gin.Context) {
	utils.SuccessResponse(c, http.StatusOK, "Sehat Bro", "Sehat")
}
