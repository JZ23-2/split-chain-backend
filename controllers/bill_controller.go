package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/services"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

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

// AssignParticipantToItem godoc
//
//	@Summary		Assign Participant To Item
//	@Description	Assign Participant To Item
//	@Tags			Bill
//	@Accept			json
//	@Produce		json
//
// @Param        bill  body      dtos.AssignParticipantsRequest  true  "Participants"
//
// @Success      200   {object}  dtos.AssignedParticipantResponse
//
//	@Failure		400			{object}	map[string]string
//	@Failure		500			{object}	map[string]string
//	@Router			/bills/assign-participants [post]
func AssignParticipantController(c *gin.Context) {
	var req dtos.AssignParticipantsRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid request")
		return
	}

	resp, err := services.AssignParticipantsToItem(req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Assign participant success", resp)
}

// GetBillsByParticipantHandler godoc
// @Summary      Get bills by participant ID
// @Description  Get bills by participant ID
// @Tags         Bill
// @Param        participantId path string true "Participant ID"
// @Produce      json
// @Success      200 {array} dtos.ParticipantBillResponse
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Router       /bills/by-participant/{participantId} [get]
func GetBillsByParticipantController(c *gin.Context) {
	participantID := c.Param("participantId")
	if participantID == "" {
		utils.FailedResponse(c, http.StatusBadRequest, "participantId is required")
		return
	}

	bills, err := services.GetBillsByParticipantID(participantID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Bill Fetched", bills)
}

// GetBillsByBillIdHandler godoc
// @Summary      Get bills by Bill ID
// @Description  Get bills by Bill ID
// @Tags         Bill
// @Param        billId path string true "Bill ID"
// @Produce      json
// @Success      200 {array} dtos.ParticipantBillResponse
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Router       /bills/by-billId/{billId} [get]
func GetBillByBillIDHandler(c *gin.Context) {
	billID := c.Param("billId")
	if billID == "" {
		utils.FailedResponse(c, http.StatusBadRequest, "billId is required")
		return
	}

	bill, err := services.GetBillByBIllID(billID)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Bill fetched successfully", bill)

}

// DeleteBillByBillIdHandler godoc
// @Summary      Delete bill by Bill ID
// @Description  Delete bill by Bill ID
// @Tags         Bill
// @Param        billId path string true "Bill ID"
// @Produce      json
// @Success      200 {object} map[string]string "Successfully deleted bill"
// @Failure      400 {object} map[string]string "Bad Request"
// @Failure      404 {object} map[string]string "Not Found"
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Router       /bills/delete-bill/{billId} [delete]
func DeleteBillByIDController(c *gin.Context) {
	billId := c.Param("billId")

	message, statusCode, err := services.DeleteBillByIDService(billId)

	if err != nil {
		utils.FailedResponse(c, statusCode, message)
		return
	}

	utils.SuccessResponse(c, statusCode, message, nil)

}

// UpdateBillByBillIdHandler godoc
// @Summary      Update Bill
// @Description  Update Bill
// @Tags         Bill
// @Param        bill body dtos.UpdateBillRequest true "Bill Data"
// @Produce      json
// @Success      200 {object} dtos.UpdateBillResponse
// @Failure      500 {object} map[string]string "Internal Server Error"
// @Router       /bills/update-bill [patch]
func UpdateBillController(c *gin.Context) {
	var req dtos.UpdateBillRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	resp, err := services.UpdateBillService(req)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "Failed to update bill")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Update bill success", resp)
}
