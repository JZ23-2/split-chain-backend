package services

import (
	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/models"
)

func GetUnreadInboxesByUserID(userID string) ([]models.Inbox, error) {
	var inboxes []models.Inbox
	if err := database.DB.
		Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Find(&inboxes).Error; err != nil {
		return nil, err
	}
	return inboxes, nil
}

func MarkInboxAsRead(inboxID string) error {
	return database.DB.
		Model(&models.Inbox{}).
		Where("id = ?", inboxID).
		Update("is_read", true).Error
}
