package models

type User struct {
	WalletAddress string `gorm:"primaryKey;type:varchar(255)" json:"wallet_address"`
}
