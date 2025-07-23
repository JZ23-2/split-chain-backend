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
//
//	@Summary	Accept friend request
//
// Description Accept a friend request
//
//	@Tags		Friend
//	@Accept		json
//	@Produce	json
//	@Param		friend	body		dtos.AcceptFriendRequest	true	"Friend Info"
//	@Success	201		{object}	dtos.AcceptFriendResponse
//	@Failure	400		"Invalid Request"
//	@Failure	404		"User or Friend Not Found"
//	@Failure	409		"Relationship Already Exists"
//	@Failure	500		"Internal Server Error"
//	@Router		/friends/accept [post]
func AcceptFriendRequest(c *gin.Context) {
	var req dtos.AcceptFriendRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	var pending models.PendingFriendRequest
	if err := database.DB.First(&pending, "id = ?", req.ID).Error; err != nil {
		utils.FailedResponse(c, http.StatusNotFound, "Not Found")
		return
	}

	friend1 := models.Friend{
		UserWalletAddress:   pending.FriendWalletAddress,
		FriendWalletAddress: pending.UserWalletAddress,
	}

	friend2 := models.Friend{
		UserWalletAddress:   pending.UserWalletAddress,
		FriendWalletAddress: pending.FriendWalletAddress,
	}

	if err := database.DB.Create(&friend1).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "Failed to create friends")
		return
	}

	if err := database.DB.Create(&friend2).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "Failed to create friends")
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

	utils.SuccessResponse(c, http.StatusOK, "Friend Request Accepted", gin.H{
		"friend_1": friend1Resp,
		"friend_2": friend2Resp,
	})
}

// DeclineFriendRequest godoc
//
//	@Summary	Decline friend request
//
// Description Decline a friend request
//
//	@Tags		Friend
//	@Accept		json
//	@Produce	json
//	@Param		friend	body		dtos.DeclineFriendRequest	true	"Friend Info"
//	@Success	201		{object}	dtos.DeclineFriendResponse
//	@Failure	400		"Invalid Request"
//	@Failure	404		"User or Friend Not Found"
//	@Failure	409		"Relationship Already Exists"
//	@Failure	500		"Internal Server Error"
//	@Router		/friends/decline [post]
func DeclineFriendRequest(c *gin.Context) {
	var req dtos.DeclineFriendRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	var pending models.PendingFriendRequest
	if err := database.DB.First(&pending, "id = ?", req.ID).Error; err != nil {
		utils.FailedResponse(c, http.StatusNotFound, "Not Found")
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

// AddFriend godoc
//
//	@Summary	Create friend request
//
// Description Create a new friend request
//
//	@Tags		Friend
//	@Accept		json
//	@Produce	json
//	@Param		friend	body		dtos.AddFriendRequest	true	"Friend Info"
//	@Success	201		{object}	dtos.AddFriendResponse
//	@Failure	400		"Invalid Request"
//	@Failure	404		"User or Friend Not Found"
//	@Failure	409		"Relationship Already Exists"
//	@Failure	500		"Internal Server Error"
//	@Router		/friends/add [post]
func AddFriend(c *gin.Context) {
	var req dtos.AddFriendRequest
	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	var user, friend models.User
	if err := database.DB.First(&user, "wallet_address = ?", req.UserWalletAddress).Error; err != nil {
		utils.FailedResponse(c, http.StatusNotFound, "User Not Found")
		return
	}

	if err := database.DB.First(&friend, "wallet_address = ?", req.FriendWalletAddress).Error; err != nil {
		utils.FailedResponse(c, http.StatusNotFound, "Friend Not Found")
		return
	}

	var existing models.PendingFriendRequest
	if err := database.DB.Where("user_wallet_address = ? AND friend_wallet_address = ? AND status = ?", req.UserWalletAddress, req.FriendWalletAddress, "Declined").First(&existing).Error; err == nil {
		existing.Status = "Pending"
		if err := database.DB.Save(&existing).Error; err != nil {
			utils.FailedResponse(c, http.StatusInternalServerError, "Failed")
			return
		}
		utils.SuccessResponse(c, http.StatusOK, "Successfully Added Friend Request", existing)
		return
	}

	if err := database.DB.Where("user_wallet_address = ? AND friend_wallet_address = ?", req.UserWalletAddress, req.FriendWalletAddress).First(&existing).Error; err == nil {
		utils.FailedResponse(c, http.StatusConflict, "Friend Request Already Send")
		return
	}

	var friendCheck models.Friend
	if err := database.DB.Where("user_wallet_address = ? AND friend_wallet_address = ?", req.UserWalletAddress, req.FriendWalletAddress).First(&friendCheck).Error; err == nil {
		utils.FailedResponse(c, http.StatusConflict, "Already Friend")
		return
	}

	newRequest := models.PendingFriendRequest{
		UserWalletAddress:   req.UserWalletAddress,
		FriendWalletAddress: req.FriendWalletAddress,
	}

	if err := database.DB.Create(&newRequest).Error; err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "Failed to create friend request")
		return
	}

	response := dtos.AddFriendResponse{
		ID:                  newRequest.ID,
		UserWalletAddress:   newRequest.UserWalletAddress,
		FriendWalletAddress: newRequest.FriendWalletAddress,
		Status:              "Pending",
	}

	utils.SuccessResponse(c, http.StatusOK, "Successfully Added Friend Request", response)
}

// TODO: Fetch Friend for user
// Routes: /friend/{user_id}
// Return: Friend array

// TODO: Add alias for friend
// Routes: /friend/alias/{user_id}/{friend_id}
// Description: Jackson mau update nama VK di tempat friend nya jadi apa
