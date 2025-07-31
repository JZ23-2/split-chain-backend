package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bill struct {
	BillID              string    `gorm:"primaryKey;type:varchar(255)" json:"billId"`
	BillTitle           string    `gorm:"type:varchar(100)" json:"billTitle"`
	StoreName           string    `gorm:"type:varchar(200)" json:"storeName"`
	TotalAmount         int       `gorm:"type:int(10)" json:"totalAmount"`
	TotalAmountAfterTax int       `gorm:"type:int(10)" json:"totalAmountAfterTax"`
	CreatorID           string    `gorm:"type:varchar(255)" json:"creatorId"`
	Creator             User      `gorm:"foreignKey:CreatorID" json:"creator"`
	CreatedAt           time.Time `json:"createdAt"`
	Tax                 float32   `json:"tax"`
	Service             float32   `json:"service"`
	Items               []Item    `gorm:"foreignKey:BillID;references:BillID" json:"items"`
}

func (b *Bill) BeforeCreate(tx *gorm.DB) (err error) {
	b.BillID = uuid.NewString()
	fmt.Println(b.BillID)
	return
}
