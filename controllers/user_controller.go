package controllers

import (
	"errors"
	"net/http"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterUser godoc
//	@Summary		Register a new user
//	@Description	Save wallet address to database
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User	true	"User info"
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{object}	map[string]string
//	@Failure		409		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/users [post]
func RegisterUser(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid Input")
		return
	}

	var existing models.User
	err := database.DB.Where("wallet_address = ?", input.WalletAddress).First(&existing).Error
	if err == nil {
		utils.FailedResponse(c, http.StatusConflict, "Wallet already registered")
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		utils.FailedResponse(c, http.StatusInternalServerError, "Database error")
		return
	}

	if err := database.DB.Create(&input).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "Failed to register user")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "User registered successfully", input)
}
