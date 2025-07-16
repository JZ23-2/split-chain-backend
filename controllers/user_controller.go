package controllers

import (
	"errors"
	"net/http"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Save wallet address to database
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body models.User true "User info"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /register [post]
func RegisterUser(c *gin.Context) {
	var input models.User

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if input.Wallet == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wallet is required"})
		return
	}

	input.UserID = uuid.New().String()

	var existing models.User
	err := database.DB.Where("wallet = ?", input.Wallet).First(&existing).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Wallet already registered"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if err := database.DB.Create(&input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User registered successfully",
		"user":    input,
	})
}
