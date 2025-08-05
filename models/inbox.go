package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Inbox struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	UserID    string    `gorm:"index" json:"userId"`
	Message   string    `gorm:"type:text" json:"message"`
	Link      string    `gorm:"type:varchar(255)" json:"link"`
	Type      string    `gorm:"type:varchar(100)" json:"type"`
	CreatedAt time.Time `json:"createdAt"`
	IsRead    bool      `gorm:"type:boolean" json:"isRead"`
}

func (i *Inbox) BeforeCreate(tx *gorm.DB) (err error) {
	i.ID = uuid.NewString()
	i.CreatedAt = time.Now()
	return
}
