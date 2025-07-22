package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/services"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

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