package handlers

import (
	"fmt"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/JZ23-2/splitbill-backend/websockets"
	"github.com/gin-gonic/gin"
)

func SendInboxToUser(userID, message, link, notifType string) {
	inbox := models.Inbox{
		UserID:  userID,
		Message: message,
		Link:    link,
		Type:    notifType,
	}

	if err := database.DB.Create(&inbox).Error; err != nil {
		fmt.Println("Gagal simpan inbox:", err)
		return
	}

	err := websockets.SendToUser(userID, gin.H{
		"type": "inbox_notification",
		"data": inbox,
	})

	if err != nil {
		fmt.Println("Gagal kirim via websocket:", err)
	}
}
