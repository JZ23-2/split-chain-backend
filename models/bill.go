package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bill struct {
	BillID       string        `gorm:"primaryKey;type:varchar(255)" json:"billId"`
	BillTitle    string        `gorm:"type:varchar(100)" json:"billTitle"`
	TotalAmount  int           `gorm:"type:int(10)" json:"totalAmount"`
	CreatorID    string        `gorm:"type:varchar(255)" json:"creatorId"`
	Creator      User          `gorm:"foreignKey:CreatorID" json:"creator"`
	CreatedAt    time.Time     `json:"createdAt"`
	Participants []Participant `gorm:"foreignKey:BillID" json:"participants"`
}

func (b *Bill) BeforeCreate(tx *gorm.DB) (err error) {
	b.BillID = uuid.NewString()
	return
}
