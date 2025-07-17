package dtos

import "github.com/JZ23-2/splitbill-backend/models"

type ParticipantDetailResponse struct {
	BillID        string        `json:"billId"`
	BillTitle     string        `json:"billTitle"`
	CreatorID     string        `json:"creatorId"`
	ParticipantID string        `json:"participantId"`
	Items         []models.Item `json:"items"`
	TotalOwed     int           `json:"totalOwed"`
}
