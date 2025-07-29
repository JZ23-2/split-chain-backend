package dtos

type AssignParticipantsRequest struct {
	BillID       string                           `json:"billId"`
	Participants []AssignParticipantDetailRequest `json:"participants"`
}

type AssignParticipantDetailRequest struct {
	ParticipantID string                         `json:"participantId"`
	IsPaid        bool                           `json:"isPaid"`
	Items         []AssignParticipantItemRequest `json:"items"`
}

type AssignParticipantItemRequest struct {
	ItemID string `json:"itemId"`
}

type AssignParticipantsResponse struct {
	BillID       string                            `json:"billId"`
	Participants []AssignParticipantDetailResponse `json:"participants"`
}

type AssignParticipantDetailResponse struct {
	ParticipantID string                            `json:"participantId"`
	IsPaid        bool                              `json:"isPaid"`
	AmountOwed    int                               `json:"amountOwed"`
	Items         []AssignParticipantItemWithAmount `json:"items"`
}

type AssignParticipantItemWithAmount struct {
	ItemID string `json:"itemId"`
	Amount int    `json:"amount"`
}
