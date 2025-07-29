package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/services"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

// AssignParticipantsToBill assigns participants to a bill and splits item prices equally among them.
// @Summary Assign participants to a bill
// @Description Assign participants to a bill and split item prices equally per participant
// @Tags Bill
// @Accept json
// @Produce json
// @Param request body dtos.AssignParticipantsRequest true "Assign Participants Request"
// @Success 201 {object} dtos.AssignParticipantsResponse "Structured Assign Participants Result"
// @Failure		400		"Invalid input"
//
//	@Failure		500		"Internal error"
//
// @Router /bills/assign-participants [post]
func AssignParticipantsController(c *gin.Context) {
	var req dtos.AssignParticipantsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid request: "+err.Error())
		return
	}

	resp, err := services.AssignParticipantToBill(req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Participants assigned successfully", resp)
}

// CreateBillWithoutParticipant godoc
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
func CreateBillWithoutParticipantController(c *gin.Context) {
	var req dtos.CreateBillWithoutParticipantRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "invalid request: "+err.Error())
		return
	}

	resp, err := services.CreateBillWithoutParticipant(req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err.Error())
		return
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
func GetBillByCreatorController(c *gin.Context) {
	creatorId := c.Query("creatorId")
	billId := c.Query("billId")

	resp, err := services.GetBillsByCreator(creatorId, billId)
	if err != nil {
		if err.Error() == "creatorId is required" {
			utils.FailedResponse(c, http.StatusBadRequest, err.Error())
			return
		}
		utils.FailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Bill fetched", resp)
}
