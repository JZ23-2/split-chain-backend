package dtos

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
	Name     string `json:"name" example:"Steak"`
	Quantity int    `json:"quantity" example:"2"`
	Price    int    `json:"price" example:"40000"`
}

type CreateBillWithoutParticipantRequest struct {
	BillTitle   string                                    `json:"billTitle" example:"Makan orang hitam bersama"`
	StoreName   string                                    `json:"storeName" example:"East Repair Inc."`
	Date        string                                    `json:"date" example:"2019-11-02"`
	Tax         float32                                   `json:"tax" example:"9.06"`
	Service     float32                                   `json:"service" example:"0.0"`
	TotalAmount int                                       `json:"totalAmount" example:"15406"`
	CreatorID   string                                    `json:"creatorId" example:"user123"`
	Items       []CreateBillWithoutParticipantItemRequest `json:"items"`
}

type CreateBillWithoutParticipantItemResponse struct {
	ItemID        string `json:"itemId"`
	Name          string `json:"name"`
	Quantity      int    `json:"quantity"`
	UnitPrice     int    `json:"price"`
	PriceAfterTax int    `json:"priceAfterTax"`
}

type CreateBillWithoutParticipantResponse struct {
	BillID              string                                     `json:"billId"`
	BillTitle           string                                     `json:"billTitle"`
	StoreName           string                                     `json:"storeName"`
	Date                string                                     `json:"date"`
	Tax                 float32                                    `json:"tax"`
	Service             float32                                    `json:"service"`
	TotalAmount         int                                        `json:"totalAmount"`
	TotalAmountAfterTax int                                        `json:"totalAmountAfterTax"`
	CreatedAt           string                                     `json:"createdAt"`
	CreatorID           string                                     `json:"creatorId"`
	Items               []CreateBillWithoutParticipantItemResponse `json:"items"`
}

type GetBillByCreatorItemResponse struct {
	ItemID   string `json:"itemId"`
	Name     string `json:"name"`
	Price    int    `json:"price"`
	Quantity int    `json:"quantity"`
}

type GetBillByCreatorResponse struct {
	BillID      string                         `json:"billId"`
	BillTitle   string                         `json:"billTitle"`
	TotalAmount int                            `json:"totalAmount"`
	Tax         float32                        `json:"tax"`
	Service     float32                        `json:"service"`
	CreatedAt   string                         `json:"createdAt"`
	Items       []GetBillByCreatorItemResponse `json:"items"`
}
