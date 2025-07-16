package models

type User struct {
	UserID string `gorm:"primaryKey;type:varchar(255)" json:"userId"`
	Wallet string `gorm:"unique;type:varchar(50)" json:"wallet"`
}
