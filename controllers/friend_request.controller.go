package controllers

import (
	"net/http"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/models"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
)

// AddFriend godoc
// @Summary Create friend request
// Description Create a new friend request
// @Tags Pending Friend Request
// @Accept json
// @Produce json
// @Param friend body dtos.AddFriendRequest true "Friend Info"
// @Success 201 {object} dtos.AddFriendResponse
// @Failure 400 "Invalid Request"
// @Failure 404 "User or Friend Not Found"
// @Failure 409 "Relationship Already Exists"
// @Failure 500 "Internal Server Error"
// @Router /pending-friend-request/send-request [post]
func SendFriendRequest(c *gin.Context) {
	var req dtos.AddFriendRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	var user, friend models.User
	if err := database.DB.First(&user, "wallet_address = ?", req.UserWalletAddress).Error; err != nil {
		utils.FailedResponse(c, http.StatusNotFound, "user Not Found")
		return
	}
	if err := database.DB.First(&friend, "wallet_address = ?", req.FriendWalletAddress).Error; err != nil {
		utils.FailedResponse(c, http.StatusNotFound, "friend Not Found")
		return
	}

	var existing models.PendingFriendRequest
	if err := database.DB.Where("user_wallet_address = ? AND friend_wallet_address = ?", req.UserWalletAddress, req.FriendWalletAddress).First(&existing).Error; err == nil {
		utils.FailedResponse(c, http.StatusConflict, "friend request already send")
		return
	}

	var friendCheck models.Friend
	if err := database.DB.Where("user_wallet_address = ? AND friend_wallet_address = ?", req.UserWalletAddress, req.FriendWalletAddress).First(&friendCheck).Error; err == nil {
		utils.FailedResponse(c, http.StatusConflict, "already friend")
		return
	}

	newRequest := models.PendingFriendRequest{
		UserWalletAddress:   req.UserWalletAddress,
		FriendWalletAddress: req.FriendWalletAddress,
	}

	if err := database.DB.Create(&newRequest).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "failed to create friend request")
		return
	}

	response := dtos.AddFriendResponse{
		ID:                  newRequest.ID,
		UserWalletAddress:   newRequest.UserWalletAddress,
		FriendWalletAddress: newRequest.FriendWalletAddress,
		Status:              newRequest.Status,
	}

	utils.SuccessResponse(c, http.StatusOK, "saved friend success", response)
}
