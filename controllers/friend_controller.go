package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

// AcceptFriendRequest godoc
// @Summary Accept friend request
// Description Accept a friend request
// @Tags Friend
// @Accept json
// @Produce json
// @Param friend body dtos.AcceptFriendRequest true "Friend Info"
// @Success 201 {object} dtos.AcceptFriendResponse
// @Failure 400 "Invalid Request"
// @Failure 404 "User or Friend Not Found"
// @Failure 409 "Relationship Already Exists"
// @Failure 500 "Internal Server Error"
// @Router /friend/accept-friend [post]
func AcceptFriendRequest(c *gin.Context) {
	var req dtos.AcceptFriendRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "invalid request")
		return
	}

	var pending models.PendingFriendRequest
	if err := database.DB.First(&pending, "id = ?", req.ID).Error; err != nil {
		utils.FailedResponse(c, http.StatusNotFound, "pending friend request not found")
		return
	}

	friend1 := models.Friend{
		UserWalletAddress:   req.UserWalletAddress,
		FriendWalletAddress: req.FriendWalletAddress,
	}

	friend2 := models.Friend{
		UserWalletAddress:   req.FriendWalletAddress,
		FriendWalletAddress: req.UserWalletAddress,
	}

	if err := database.DB.Create(&friend1).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "failed to create friend")
		return
	}

	if err := database.DB.Create(&friend2).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "failed to create friend")
		return
	}

	friend1Resp := dtos.AcceptFriendResponse{
		ID:                  friend1.ID,
		UserWalletAddress:   friend1.UserWalletAddress,
		FriendWalletAddress: friend1.FriendWalletAddress,
		Nickname:            &friend1.Nickname,
	}

	friend2Resp := dtos.AcceptFriendResponse{
		ID:                  friend2.ID,
		UserWalletAddress:   friend2.UserWalletAddress,
		FriendWalletAddress: friend2.FriendWalletAddress,
		Nickname:            &friend2.Nickname,
	}

	database.DB.Delete(&pending)

	utils.SuccessResponse(c, http.StatusOK, "friend request accepted", gin.H{
		"friend_1": friend1Resp,
		"friend_2": friend2Resp,
	})
}

// DeclineFriendRequest godoc
// @Summary Decline friend request
// Description Decline a friend request
// @Tags Friend
// @Accept json
// @Produce json
// @Param friend body dtos.DeclineFriendRequest true "Friend Info"
// @Success 201 {object} dtos.DeclineFriendResponse
// @Failure 400 "Invalid Request"
// @Failure 404 "User or Friend Not Found"
// @Failure 409 "Relationship Already Exists"
// @Failure 500 "Internal Server Error"
// @Router /friend/decline-friend [post]
func DeclineFriendRequest(c *gin.Context) {
	var req dtos.DeclineFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	var pending models.PendingFriendRequest
	if err := database.DB.Where(
		"user_wallet_address = ? AND friend_wallet_address = ? AND status = ?",
		req.UserWalletAddress, req.FriendWalletAddress, "Pending",
	).First(&pending).Error; err != nil {
		utils.FailedResponse(c, http.StatusNotFound, "Pending friend request not found")
		return
	}

	pending.Status = "Declined"
	if err := database.DB.Save(&pending).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "Failed to decline friend request")
		return
	}

	response := dtos.DeclineFriendResponse{
		ID:                  pending.ID,
		UserWalletAddress:   pending.UserWalletAddress,
		FriendWalletAddress: pending.FriendWalletAddress,
		Status:              pending.Status,
	}

	utils.SuccessResponse(c, http.StatusOK, "Friend request declined", response)
}
