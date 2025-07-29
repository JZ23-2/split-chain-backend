package controllers

import (
	"fmt"
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
//
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
				ParticipantID: &p.ParticipantID,
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

// CreateSimpleBill godoc
// @Summary      Create bill (no participants)
// @Description  Save a bill with items, tax, and service, without splitting between participants
// @Tags         Bill
// @Accept       json
// @Produce      json
// @Param        bill  body      dtos.CreateBillWithoutParticipantRequest  true  "Bill Data without participant"
// @Success      201   {object}  dtos.CreateBillWithoutParticipantResponse
// @Failure      400   {object}  map[string]string
// @Failure      500   {object}  map[string]string
// @Router       /bills/bill-without-participant [post]
func CreateBillWithoutParticipant(c *gin.Context) {
	var req dtos.CreateBillWithoutParticipantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	bill := models.Bill{
		BillTitle:   req.StoreName,
		TotalAmount: req.TotalAmount,
		Tax:         req.Tax,
		Service:     req.Service,
		CreatorID:   req.CreatorID,
		CreatedAt:   time.Now(),
	}

	if err := database.DB.Create(&bill).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "failed to create bill: "+err.Error())
		return
	}

	var itemResponses []dtos.CreateBillWithoutParticipantItemResponse

	fmt.Println("Bill ID: ", bill.BillID)
	for _, item := range req.Items {
		itemID := uuid.NewString()
		newItem := models.Item{
			ItemID:        itemID,
			BillID:        bill.BillID,
			ParticipantID: nil,
			Name:          item.Name,
			Price:         item.UnitPrice,
			Quantity:      item.Quantity,
		}

		if err := database.DB.Create(&newItem).Error; err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, "failed to create item: "+err.Error())
			return
		}

		itemResponses = append(itemResponses, dtos.CreateBillWithoutParticipantItemResponse{
			ItemID:        itemID,
			Name:          item.Name,
			Quantity:      item.Quantity,
			UnitPrice:     item.UnitPrice,
			PriceAfterTax: item.PriceAfterTax,
			PriceInHBAR:   item.PriceInHBAR,
		})
	}

	resp := dtos.CreateBillWithoutParticipantResponse{
		BillID:      bill.BillID,
		StoreName:   req.StoreName,
		Date:        req.Date,
		Tax:         req.Tax,
		Service:     req.Service,
		TotalAmount: req.TotalAmount,
		CreatedAt:   bill.CreatedAt.Format(time.RFC3339),
		CreatorID:   req.CreatorID,
		Items:       itemResponses,
	}

	utils.SuccessResponse(c, http.StatusCreated, "Bill without participant created successfully", resp)
}

// GetBillByCreator godoc
//
//	@Summary		Get bills by creator
//	@Description	Get all bills created by a specific creator, optionally filter by billId
//	@Tags			Bill
//	@Accept			json
//	@Produce		json
//	@Param			creatorId	query		string	true	"Creator ID"
//	@Param			billId		query		string	false	"Bill ID (optional filter)"
//	@Success		200			{object}	dtos.SuccessResponse{data=[]dtos.GetBillByCreatorResponse}
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/bills/by-creator [get]
func GetBillByCreator(c *gin.Context) {
	creatorId := c.Query("creatorId")
	billId := c.Query("billId")

	if creatorId == "" {
		utils.FailedResponse(c, http.StatusBadRequest, "Creator Id is required")
		return
	}

	var bills []models.Bill
	query := database.DB.Preload("Items").Where("creator_id = ?", creatorId)

	if billId != "" {
		query = query.Where("bill_id = ?", billId)
	}

	if err := query.Find(&bills).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "Failed to fetch bills")
		return
	}

	var response []dtos.GetBillByCreatorResponse
	for _, bill := range bills {
		var itemResponses []dtos.GetBillByCreatorItemResponse
		for _, item := range bill.Items {
			itemResponses = append(itemResponses, dtos.GetBillByCreatorItemResponse{
				ItemID:   item.ItemID,
				Name:     item.Name,
				Price:    item.Price,
				Quantity: item.Quantity,
			})
		}

		response = append(response, dtos.GetBillByCreatorResponse{
			BillID:      bill.BillID,
			BillTitle:   bill.BillTitle,
			TotalAmount: bill.TotalAmount,
			Tax:         bill.Tax,
			Service:     bill.Service,
			CreatedAt:   bill.CreatedAt.Format(time.RFC3339),
			Items:       itemResponses,
		})

	}
	utils.SuccessResponse(c, http.StatusOK, "Bill Fetched", response)
}
