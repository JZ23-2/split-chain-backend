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
