package models

type Participant struct {
	ParticipantID string `gorm:"primaryKey;type:varchar(255)" json:"participantId"`
	BillID        string `gorm:"primaryKey;type:varchar(255)" json:"billId"`
	AmountOwed    int    `gorm:"type:int(10)" json:"amountOwed"`
	IsPaid        bool   `json:"isPaid"`

	ParticipantItem []ParticipantItem `gorm:"foreignKey:ParticipantID;references:ParticipantID" json:"participants"`
	User            User              `gorm:"foreignKey:ParticipantID"`
	Bill            Bill              `gorm:"foreignKey:BillID"`
}
