package models

type User struct {
	UserID string `gorm:"primaryKey" json:"userId"`
	Wallet string `gorm:"unique" json:"wallet"`
}
