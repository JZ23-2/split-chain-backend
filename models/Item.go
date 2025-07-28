package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ItemID        string  `gorm:"primaryKey;type:varchar(255)" json:"itemId"`
	ParticipantID *string `gorm:"type:varchar(255);default:null" json:"participantId"`
	BillID        string  `gorm:"type:varchar(255)" json:"billId"`
	Name          string  `gorm:"type:varchar(255)" json:"name"`
	Quantity      int     `gorm:"type:int(10)" json:"quantity"`
	Price         int     `gorm:"type:int(10)" json:"price"`
}

func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	i.ItemID = uuid.NewString()
	return
}
