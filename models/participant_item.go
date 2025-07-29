package models

type ParticipantItem struct {
	ParticipantID string `gorm:"type:varchar(255)" json:"participantId"`
	ItemID        string `gorm:"type:varchar(255)" json:"itemId"`
	Amount        int    `gorm:"type:int(10)" json:"amount"`
}
