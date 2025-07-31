package dtos

type ParticipantBillResponse struct {
	BillID      string                    `json:"billId"`
	BillTitle   string                    `json:"billTitle"`
	TotalAmount int                       `json:"totalAmount"`
	CreatorID   string                    `json:"creatorId"`
	Tax         float32                   `json:"tax"`
	Service     float32                   `json:"service"`
	Item        []ParticipantItemResponse `json:"items"`
	Participant []ParticipantListResponse `json:"participants"`
}

type ParticipantListResponse struct {
	ParticipantID string `json:"participantId"`
	AmountOwed    int    `json:"amountOwed"`
	IsPaid        bool   `json:"isPaid"`
}

type ParticipantItemResponse struct {
	ItemID        string `gorm:"primaryKey;type:varchar(255)" json:"itemId"`
	Name          string `gorm:"type:varchar(255)" json:"name"`
	Quantity      int    `gorm:"type:int(10)" json:"quantity"`
	Price         int    `gorm:"type:int(10)" json:"price"`
	PriceAfterTax int    `gorm:"type:int(10)" json:"priceAfterTax"`
}
