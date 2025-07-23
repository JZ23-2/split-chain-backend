package controllers

import (
	"net/http"
	"time"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// CreateBill godoc
//	@Summary		Create a new bill
//	@Description	Create a bill with participants and their items
//	@Tags			Bill
//	@Accept			json
//	@Produce		json
//	@Param			bill	body		dtos.CreateBillRequest	true	"Bill Info"
//	@Success		201		{object}	dtos.CreateBillResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		500		{object}	map[string]string
//	@Router			/bills [post]
func CreateBill(c *gin.Context) {
	var req dtos.CreateBillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	bill := models.Bill{
		BillTitle:   req.BillTitle,
		TotalAmount: req.TotalAmount,
		CreatorID:   req.CreatorID,
		CreatedAt:   time.Now(),
	}

	bill.BillID = uuid.NewString()

	for _, p := range req.Participants {
		participant := models.Participant{
			ParticipantID: p.ParticipantID,
			BillID:        bill.BillID,
			AmountOwed:    p.AmountOwed,
			IsPaid:        p.IsPaid,
		}

		for _, item := range p.Items {
			participant.Items = append(participant.Items, models.Item{
				ItemID:        uuid.NewString(),
				ParticipantID: p.ParticipantID,
				Name:          item.Name,
				Price:         item.Price,
			})
		}

		bill.Participants = append(bill.Participants, participant)
	}

	if err := database.DB.Create(&bill).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "failed to create bill")
		return
	}

	resp := dtos.CreateBillResponse{
		BillID:      bill.BillID,
		BillTitle:   bill.BillTitle,
		TotalAmount: bill.TotalAmount,
		CreatorID:   bill.CreatorID,
		CreatedAt:   bill.CreatedAt.Format(time.RFC3339),
	}

	for _, p := range bill.Participants {
		participantResp := dtos.CreateBillParticipantResponse{
			ParticipantID: p.ParticipantID,
			AmountOwed:    p.AmountOwed,
			IsPaid:        p.IsPaid,
		}

		for _, item := range p.Items {
			participantResp.Items = append(participantResp.Items, dtos.CreateBillItemResponse{
				ItemID: item.ItemID,
				Name:   item.Name,
				Price:  item.Price,
			})
		}

		resp.Participants = append(resp.Participants, participantResp)
	}

	utils.SuccessResponse(c, http.StatusCreated, "Bill created successfully", resp)
}
