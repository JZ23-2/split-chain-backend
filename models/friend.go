package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friend struct {
	ID                  string `gorm:"PrimaryKey;type:varchar(255)" json:"id"`
	UserWalletAddress   string `gorm:"type:varchar(255)" json:"user_wallet_address"`
	FriendWalletAddress string `gorm:"type:varchar(255)" json:"friend_wallet_address"`
	Nickname            string `gorm:"type:varchar(100);default:NULL" json:"nickname"`
	User                User   `gorm:"foreignKey:UserWalletAddress;references:WalletAddress"`
	Friend              User   `gorm:"foreignKey:FriendWalletAddress;references:WalletAddress"`
}

func (b *Friend) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.NewString()
	return
}
