package dtos

type GetParticipantDetailRequest struct {
	BillID        string `json:"billId" binding:"required"`
	ParticipantID string `json:"participantId" binding:"required"`
}
