package models

type Participant struct {
	ParticipantID string `gorm:"primaryKey;type:varchar(255)" json:"participantId"`
	User          User   `gorm:"foreignKey:ParticipantID"`
	BillID        string `gorm:"type:varchar(255)" json:"billId"`
	AmountOwed    int    `gorm:"type:int(10)" json:"amountOwed"`
	IsPaid        bool   `json:"isPaid"`
	Items         []Item `gorm:"foreignKey:ParticipantID"`
}
