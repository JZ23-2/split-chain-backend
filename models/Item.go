package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ItemID        string `gorm:"primaryKey;type:varchar(255)" json:"itemId"`
	ParticipantID string `gorm:"type:varchar(255)" json:"participantId"`
	BillID        string `gorm:"type:varchar(255)" json:"billId"`
	Name          string `gorm:"type:varchar(255)" json:"name"`
	Price         int    `gorm:"type:int(10)" json:"price"`
}

func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	i.ItemID = uuid.NewString()
	return
}
