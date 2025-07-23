package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/services"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

// ConvertPricesToHBAR godoc
//	@Summary		Convert item prices to HBAR
//	@Description	Takes a receipt JSON, converts each item's price (after tax) to HBAR using the current rate, and returns the updated receipt.
//	@Tags			Rate Conversion
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.ReceiptResponse	true	"Receipt data"
//	@Success		200		{object}	dtos.ReceiptResponse
//	@Failure		400		"Invalid request"
//	@Failure		500		"Failed to fetch HBAR rate"
//	@Router			/rate [post]
func ConvertPricesToHBAR(c *gin.Context) {
	var req dtos.ReceiptResponse

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	rate, err := services.FetchHBARRate()
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "failed to fetch HBAR rate")
		return
	}

	for i := range req.Items {
		hbar := req.Items[i].PriceAfterTax * float32(rate)
		req.Items[i].PriceInHBAR = hbar
	}

	utils.SuccessResponse(c, http.StatusOK, "success get HBAR rate", req)
}
