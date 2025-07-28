package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/services"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

// ConvertPricesToHBAR godoc
//
//	@Summary		Get HBAR Rate (1 USD)
//	@Description	Get HBAR Rate by 1 USD
//	@Tags			Rate
//	@Accept			json
//	@Produce		json
//
// @Success		200	{object}	map[string]interface{}	"success get HBAR rate"
//
//	@Failure		400		"Invalid request"
//	@Failure		500		"Failed to fetch HBAR rate"
//	@Router			/get-rate [get]
func ConvertPricesToHBAR(c *gin.Context) {

	rate, err := services.FetchHBARRate()
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "failed to fetch HBAR rate")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success get HBAR rate", rate)
}
