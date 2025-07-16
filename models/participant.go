package models

type Participant struct {
	ParticipantID string `gorm:"primaryKey;type:varchar(255)" json:"participantId"`
	BillID        string `gorm:"type:varchar(255)" json:"billId"`
	CreatorID     string `gorm:"type:varchar(255)" json:"creatorId"`
	User          User   `gorm:"foreignKey:CreatorID;;references:UserID"`
	AmountOwed    int    `gorm:"type:int(10)" json:"amountOwed"`
	IsPaid        bool   `json:"isPaid"`
	Items         []Item `gorm:"foreignKey:ParticipantID"`
}
