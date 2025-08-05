package dtos

import "time"

type CreateItemRequest struct {
	Name  string `json:"name" example:"Steak"`
	Price int    `json:"price" example:"80000"`
}

type CreateParticipantRequest struct {
	ParticipantID string              `json:"participantId" example:"user123"`
	AmountOwed    int                 `json:"amountOwed" example:"100000"`
	IsPaid        bool                `json:"isPaid" example:"true"`
	Items         []CreateItemRequest `json:"items"`
}

type CreateBillRequest struct {
	BillTitle    string                     `json:"billTitle" example:"Dinner at Cafe"`
	TotalAmount  int                        `json:"totalAmount" example:"200000"`
	CreatorID    string                     `json:"creatorId" example:"user123"`
	Participants []CreateParticipantRequest `json:"participants"`
}

type CreateBillResponse struct {
	BillID       string                          `json:"billId"`
	BillTitle    string                          `json:"billTitle"`
	TotalAmount  int                             `json:"totalAmount"`
	CreatorID    string                          `json:"creatorId"`
	CreatedAt    string                          `json:"createdAt"`
	Participants []CreateBillParticipantResponse `json:"participants"`
}

type CreateBillParticipantResponse struct {
	ParticipantID string                   `json:"participantId"`
	AmountOwed    int                      `json:"amountOwed"`
	IsPaid        bool                     `json:"isPaid"`
	Items         []CreateBillItemResponse `json:"items"`
}

type CreateBillItemResponse struct {
	ItemID string `json:"itemId"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
}

type CreateBillWithoutParticipantItemRequest struct {
	Name     string  `json:"name" example:"Steak"`
	Quantity int     `json:"quantity" example:"2"`
	Price    float64 `json:"price" example:"10.99"`
}

type CreateBillWithoutParticipantRequest struct {
	StoreName string                                    `json:"storeName" example:"East Repair Inc."`
	BillDate  string                                    `json:"billDate" example:"2019-11-02"`
	Tax       float64                                   `json:"tax" example:"9.06"`
	Service   float64                                   `json:"service" example:"0.0"`
	CreatorID string                                    `json:"creatorId" example:"user123"`
	Items     []CreateBillWithoutParticipantItemRequest `json:"items"`
}

type CreateBillWithoutParticipantItemResponse struct {
	ItemID    string  `json:"itemId"`
	Name      string  `json:"name"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"price"`
}

type CreateBillWithoutParticipantResponse struct {
	BillID    string                                     `json:"billId"`
	StoreName string                                     `json:"storeName"`
	BillDate  time.Time                                  `json:"billDate"`
	Tax       float64                                    `json:"tax"`
	CreatedAt string                                     `json:"createdAt"`
	CreatorID string                                     `json:"creatorId"`
	Items     []CreateBillWithoutParticipantItemResponse `json:"items"`
}

type GetBillByCreatorItemResponse struct {
	ItemID       string                                `json:"itemId"`
	Name         string                                `json:"name"`
	Price        float64                               `json:"price"`
	Quantity     int                                   `json:"quantity"`
	Participants []GetBillByCreatorParticipantResponse `json:"participants"`
}

type GetBillByCreatorParticipantResponse struct {
	ParticipantID string  `json:"participantId"`
	AmountOwed    float64 `json:"amountOwed"`
	IsPaid        string    `json:"isPaid"`
}

type GetBillByCreatorResponse struct {
	BillID       string                                `json:"billId"`
	StoreName    string                                `json:"storeName"`
	Tax          float64                               `json:"tax"`
	CreatedAt    string                                `json:"createdAt"`
	BillDate     time.Time                             `json:"billDate"`
	Items        []GetBillByCreatorItemResponse        `json:"items"`
	Participants []GetBillByCreatorParticipantResponse `json:"participants"`
}
