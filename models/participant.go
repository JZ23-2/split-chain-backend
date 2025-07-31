package models

type Participant struct {
	ParticipantID string `gorm:"primaryKey;type:varchar(255)" json:"participantId"`
	ItemID        string `gorm:"primaryKey;type:varchar(255)" json:"itemId"`
	AmountOwed    int    `gorm:"type:int(10)" json:"amountOwed"`
	IsPaid        bool   `json:"isPaid"`

	User User `gorm:"foreignKey:ParticipantID"`
}
