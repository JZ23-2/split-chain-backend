package dtos

import "time"

type ParticipantBillResponse struct {
	BillID       string                    `json:"billId"`
	StoreName    string                    `json:"storeName"`
	CreatorID    string                    `json:"creatorId"`
	BillDate     time.Time                 `json:"billDate"`
	CreatedAt    string                    `json:"createdAt"`
	Tax          float64                   `json:"tax"`
	Items        []ParticipantItemResponse `json:"items"`
	Participants []ParticipantListResponse `json:"participants"`
}

type ParticipantItemResponse struct {
	ItemID       string                    `json:"itemId"`
	Name         string                    `json:"name"`
	Quantity     int                       `json:"quantity"`
	Price        float64                   `json:"price"`
	Participants []ParticipantListResponse `json:"participants"`
}

type ParticipantListResponse struct {
	ParticipantID string  `json:"participantId"`
	AmountOwed    float64 `json:"amountOwed"`
	IsPaid        string    `json:"isPaid"`
}
