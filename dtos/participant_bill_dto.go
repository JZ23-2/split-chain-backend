package dtos

import "time"

type ParticipantBillResponse struct {
	BillID       string                    `json:"billId"`
	StoreName    string                    `json:"storeName"`
	CreatorID    string                    `json:"creatorId"`
	BillDate     time.Time                 `json:"billDate"`
	CreatedAt    string                    `json:"createdAt"`
	Tax          int                       `json:"tax"`
	DisplayTax   string                    `json:"displayTax"`
	Items        []ParticipantItemResponse `json:"items"`
	Participants []ParticipantListResponse `json:"participants"`
}

type ParticipantItemResponse struct {
	ItemID       string                    `json:"itemId"`
	Name         string                    `json:"name"`
	Quantity     int                       `json:"quantity"`
	Price        int                       `json:"price"`
	DisplayPrice string                    `json:"displayPrice"`
	Participants []ParticipantListResponse `json:"participants"`
}

type ParticipantListResponse struct {
	ParticipantID     string `json:"participantId"`
	AmountOwed        int    `json:"amountOwed"`
	DisplayAmountOwed string `json:"displayAmountOwed"`
	IsPaid            bool   `json:"isPaid"`
}
