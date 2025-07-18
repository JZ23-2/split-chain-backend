package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/services"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

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
		utils.FailedResponse(c, http.StatusInternalServerError, "failed to process image:" + err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "result", result)
}