package dtos

import "time"

type ParticipantBillResponse struct {
	BillID      string                    `json:"billId"`
	StoreName   string                    `json:"storeName"`
	CreatorID   string                    `json:"creatorId"`
	BillDate    time.Time                 `json:"billDate"`
	CreatedAt   string                    `json:"createdAt"`
	Tax         float32                   `json:"tax"`
	Item        []ParticipantItemResponse `json:"items"`
	Participant []ParticipantListResponse `json:"participants"`
}

type ParticipantListResponse struct {
	ParticipantID string `json:"participantId"`
	AmountOwed    int    `json:"amountOwed"`
	IsPaid        bool   `json:"isPaid"`
}

type ParticipantItemResponse struct {
	ItemID   string `gorm:"primaryKey;type:varchar(255)" json:"itemId"`
	Name     string `gorm:"type:varchar(255)" json:"name"`
	Quantity int    `gorm:"type:int(10)" json:"quantity"`
	Price    int    `gorm:"type:int(10)" json:"price"`
}
