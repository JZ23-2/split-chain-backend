package models

import "time"

type Bill struct {
	BillID       string        `gorm:"primaryKey;type:varchar(255)" json:"billId"`
	Title        string        `gorm:"type:varchar(100)" json:"title"`
	TotalAmount  int           `gorm:"type:int(10)" json:"totalAmount"`
	CreatorID    string        `gorm:"type:varchar(255)" json:"creatorId"`
	Creator      User          `gorm:"foreignKey:CreatorID" json:"creator"`
	CreatedAt    time.Time     `json:"createdAt"`
	Participants []Participant `gorm:"foreignKey:BillID" json:"participants"`
}
