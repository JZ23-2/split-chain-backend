package dtos

import "time"

type UpdateBillParticipantRequest struct {
	ParticipantID string `json:"participantId" example:"u123"`
	AmountOwed    int    `json:"amountOwed" example:"15000"`
	IsPaid        bool   `json:"isPaid" example:"false"`
}

type UpdateBillItemRequest struct {
	ItemID                       string                         `json:"itemId" example:"item-001"`
	Name                         string                         `json:"name" example:"Nasi Goreng"`
	Quantity                     int                            `json:"quantity" example:"2"`
	Price                        int                            `json:"price" example:"20000"`
	UpdateBillParticipantRequest []UpdateBillParticipantRequest `json:"participants"`
}

type UpdateBillRequest struct {
	BillID                string                  `json:"billId" example:"bill-001"`
	StoreName             string                  `json:"storeName" example:"Warung Makan Bu Tini"`
	CreatorID             string                  `json:"creatorId" example:"user-123"`
	CreatedAt             time.Time               `json:"createdAt" example:"2025-07-31T15:04:05Z"`
	BillDate              time.Time               `json:"billDate" example:"2025-07-30T00:00:00Z"`
	Tax                   float32                 `json:"tax" example:"10.0"`
	UpdateBillItemRequest []UpdateBillItemRequest `json:"items"`
}

type UpdateBillParticipantResponse struct {
	ParticipantID string `json:"participantId"`
	AmountOwed    int    `json:"amountOwed"`
	IsPaid        bool   `json:"isPaid"`
}

type UpdateBillItemResponse struct {
	ItemID                        string                          `json:"itemId"`
	Name                          string                          `json:"name"`
	Quantity                      int                             `json:"quantity"`
	Price                         int                             `json:"price"`
	UpdateBillParticipantResponse []UpdateBillParticipantResponse `json:"participants"`
}

type UpdateBillResponse struct {
	BillID                 string                   `json:"billId"`
	StoreName              string                   `json:"storeName"`
	CreatorID              string                   `json:"creatorId"`
	CreatedAt              time.Time                `json:"createdAt"`
	BillDate               time.Time                `json:"billDate"`
	Tax                    float32                  `json:"tax"`
	UpdateBillItemResponse []UpdateBillItemResponse `json:"items"`
}
