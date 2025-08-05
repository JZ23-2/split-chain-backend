package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/services"
	"github.com/gin-gonic/gin"
)

func GetUnreadInboxesByUserIDHandler(c *gin.Context) {
	userID := c.Query("userId")

	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "userId is required"})
		return
	}

	inboxes, err := services.GetUnreadInboxesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch inboxes"})
		return
	}

	c.JSON(http.StatusOK, inboxes)
}

func MarkInboxAsReadHandler(c *gin.Context) {
	inboxID := c.Param("id")
	if inboxID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Inbox ID is required"})
		return
	}

	err := services.MarkInboxAsRead(inboxID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark inbox as read"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Inbox marked as read successfully"})
}
