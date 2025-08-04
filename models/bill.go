package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bill struct {
	BillID    string    `gorm:"primaryKey;type:varchar(255)" json:"billId"`
	StoreName string    `gorm:"type:varchar(200)" json:"storeName"`
	CreatorID string    `gorm:"type:varchar(255)" json:"creatorId"`
	Creator   User      `gorm:"foreignKey:CreatorID" json:"creator"`
	CreatedAt time.Time `json:"createdAt"`
	BillDate  time.Time `json:"billDate"`
	Tax       int       `json:"tax"`
	Items     []Item    `gorm:"foreignKey:BillID;references:BillID" json:"items"`
}

func (b *Bill) BeforeCreate(tx *gorm.DB) (err error) {
	b.BillID = uuid.NewString()
	fmt.Println(b.BillID)
	return
}
