package dtos

type AssignParticipantsRequest struct {
	ItemID        string   `json:"itemId"`
	ParticipantID []string `json:"participants"`
}

type AssignedParticipant struct {
	ParticipantID string `json:"participantId"`
	ItemID        string `json:"itemId"`
	AmountOwed    float64    `json:"amountOwed"`
	IsPaid        bool   `json:"isPaid"`
}

type AssignedParticipantResponse struct {
	ItemID       string                `json:"itemId"`
	ItemName     string                `json:"itemName"`
	Participants []AssignedParticipant `json:"participants"`
}
