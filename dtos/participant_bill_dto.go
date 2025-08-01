package dtos

import "time"

type ParticipantBillResponse struct {
	BillID       string                    `json:"billId"`
	StoreName    string                    `json:"storeName"`
	CreatorID    string                    `json:"creatorId"`
	BillDate     time.Time                 `json:"billDate"`
	CreatedAt    string                    `json:"createdAt"`
	Tax          float32                   `json:"tax"`
	Items        []ParticipantItemResponse `json:"items"`
	Participants []ParticipantListResponse `json:"participants"`
}

type ParticipantItemResponse struct {
	ItemID       string                    `json:"itemId"`
	Name         string                    `json:"name"`
	Quantity     int                       `json:"quantity"`
	Price        int                       `json:"price"`
	Participants []ParticipantListResponse `json:"participants"`
}

type ParticipantListResponse struct {
	ParticipantID string `json:"participantId"`
	AmountOwed    int    `json:"amountOwed"`
	IsPaid        bool   `json:"isPaid"`
}
