package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Item struct {
	ItemID   string `gorm:"primaryKey;type:varchar(255)" json:"itemId"`
	BillID   string `gorm:"type:varchar(255)" json:"billId"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Quantity int    `gorm:"type:int(10)" json:"quantity"`
	Price    int    `gorm:"type:int(10)" json:"price"`

	Participant []Participant `gorm:"foreignKey:ItemID;references:ItemID" json:"participants"`
}

func (i *Item) BeforeCreate(tx *gorm.DB) (err error) {
	i.ItemID = uuid.NewString()
	return
}
