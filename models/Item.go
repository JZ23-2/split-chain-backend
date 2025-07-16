package models

type Item struct {
	ItemID        string `gorm:"primaryKey;type:varchar(255)" json:"itemId"`
	ParticipantID string `gorm:"type:varchar(255)" json:"participantId"`
	Name          string `gorm:"type:varchar(100)" json:"name"`
	Price         int    `gorm:"type:int(10)" json:"price"`
}
