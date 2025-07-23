package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

// GetParticipantDetail godoc
//	@Summary		Get participant Bills
//	@Description	Get participant all bills
//	@Tags			Participants
//	@Accept			json
//	@Produce		json
//	@Param			participant_id	path		string	true	"Participant ID"
//	@Success		200				{object}	dtos.ParticipantDetailResponse
//	@Failure		404				{object}	map[string]string
//	@Router			/participants/{participant_id} [get]
func GetParticipantBills(c *gin.Context) {
	participantId := c.Param("participant_id")
	var participants []models.Participant

	err := database.DB.
		Where("participant_id = ?", &participantId).
		Preload("Bill").
		Preload("Items").
		Find(&participants).Error
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "DB error")
		return
	}

	if len(participants) == 0 {
		utils.FailedResponse(c, http.StatusNotFound, "No bills found for participant")
		return
	}

	resp := make([]dtos.ParticipantDetailResponse, 0, len(participants))
	for _, p := range participants {
		resp = append(resp, dtos.ParticipantDetailResponse{
			BillID:        p.BillID,
			BillTitle:     p.Bill.BillTitle,
			CreatorID:     p.Bill.CreatorID,
			ParticipantID: p.ParticipantID,
			Items:         p.Items,
			TotalOwed:     p.AmountOwed,
		})
	}

	utils.SuccessResponse(c, http.StatusOK, "Participant details retrieved successfully", resp)
}

// GetParticipantDetail godoc
//	@Summary		Get participant detail in a bill
//	@Description	Retrieve participant detail including bill and items using billId and participantId
//	@Tags			Participants
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dtos.GetParticipantDetailRequest	true	"Bill and Participant IDs"
//	@Success		200		{object}	dtos.ParticipantDetailResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		404		{object}	map[string]string
//	@Router			/participants/{participant_id}/{bill_id} [get]
func GetParticipantDetail(c *gin.Context) {
	participantId := c.Param("participant_id")
	billId := c.Param("bill_id")

	var participant models.Participant
	err := database.DB.
		Where("bill_id = ? AND participant_id = ?", billId, participantId).
		Preload("Bill").
		Preload("Items", "bill_id = ? AND participant_id = ?", billId, participantId).
		First(&participant).Error
	if err != nil {
		utils.FailedResponse(c, http.StatusNotFound, "Participant Not Found")
		return
	}

	resp := dtos.ParticipantDetailResponse{
		BillID:        participant.BillID,
		BillTitle:     participant.Bill.BillTitle,
		CreatorID:     participant.Bill.CreatorID,
		ParticipantID: participant.ParticipantID,
		Items:         participant.Items,
		TotalOwed:     participant.AmountOwed,
	}
	c.JSON(http.StatusOK, resp)
}
