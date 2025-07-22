package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PendingFriendRequest struct {
	ID                  string `gorm:"PrimaryKey;type:varchar(255)" json:"id"`
	UserWalletAddress   string `gorm:"type:varchar(255)" json:"user_wallet_address"`
	FriendWalletAddress string `gorm:"type:varchar(255)" json:"friend_wallet_address"`
	Status              string `gorm:"type:varchar(100);default:'Pending'" json:"status"`
	User                User   `gorm:"foreignKey:UserWalletAddress;references:WalletAddress"`
	Friend              User   `gorm:"foreignKey:FriendWalletAddress;references:WalletAddress"`
}

func (b *PendingFriendRequest) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.NewString()
	return
}
