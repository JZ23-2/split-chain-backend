package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/services"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

// ExtractReceipt godoc
//
//	@Summary		Extract structured receipt data from an uploaded image
//	@Description	Accepts a receipt image (PNG/JPEG) and get the detail.
//	@Tags			receipt
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file					true	"Receipt image file"
//	@Success		200		{object}	dtos.ReceiptResponse	"Structured receipt result."
//	@Failure		400		"Invalid input (missing file, bad form data)."
//	@Failure		500		"Internal error (Gemini failure, parse error, etc.)."
//	@Router			/receipt/ [post]
func ExtractReceipt(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid input")
		return
	}

	opened, err := file.Open()
	if err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Failed to read file")
		return
	}
	defer opened.Close()

	result, err := services.SendToGemini(opened)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "failed to process image:"+err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "result", result)
}
