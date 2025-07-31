package controllers

import (
	"errors"
	"net/http"

	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/services"
	"github.com/JZ23-2/splitbill-backend/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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

	friend1Resp, friend2Resp, err := services.AcceptFriendRequestService(req)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.FailedResponse(c, http.StatusNotFound, "Friend request not found")
		} else {
			utils.FailedResponse(c, http.StatusInternalServerError, "Failed to accept friend request")
		}
		return
	}

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

	response, err := services.DeclineFriendRequestService(req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.FailedResponse(c, http.StatusNotFound, "Friend request not found")
		} else {
			utils.FailedResponse(c, http.StatusInternalServerError, "Failed to decline friend request")
		}
		return
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

	result, message, statusCode, err := services.AddFriendRequestService(req)
	if err != nil {
		utils.FailedResponse(c, statusCode, message)
		return
	}

	utils.SuccessResponse(c, statusCode, message, result)
}

// FetchFriend godoc
//
//	@Summary	Fetch friend
//
// Description Fetch user friend
//
//	@Tags		Friend
//	@Accept		json
//	@Produce	json
//	@Param			user_wallet_address	path		string	true	"User Wallet Address"
//	@Success	201		{object}	dtos.FriendResponse
//	@Failure	400		"Invalid Request"
//	@Failure	404		"User or Friend Not Found"
//	@Failure	409		"Relationship Already Exists"
//	@Failure	500		"Internal Server Error"
//	@Router		/friends/{user_wallet_address} [get]
func GetFriend(c *gin.Context) {
	userWalletAddress := c.Param("user_wallet_address")

	friends, err := services.GetFriendService(userWalletAddress)
	if err != nil {
		utils.FailedResponse(c, http.StatusInternalServerError, "Failed to fetch friends")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Successfully Fetch Friend", friends)
}

// AddFriendNickname godoc
//
//	@Summary	add friend nickname
//
// Description add friend nickname
//
//	@Tags		Friend
//	@Accept		json
//	@Produce	json
//	@Param		friend	body		dtos.FriendNicknameRequest	true "Friend Info"
//	@Success	201		{object}	dtos.FriendResponse
//	@Failure	400		"Invalid Request"
//	@Failure	404		"User or Friend Not Found"
//	@Failure	409		"Relationship Already Exists"
//	@Failure	500		"Internal Server Error"
//	@Router		/friends/alias [post]
func AddFriendNickname(c *gin.Context) {
	var req dtos.FriendNicknameRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.FailedResponse(c, http.StatusBadRequest, "Invalid Request")
		return
	}

	res, err := services.AddFriendNicknameService(req)
	if err != nil {
		if err.Error() == "user or friend not found" {
			utils.FailedResponse(c, http.StatusNotFound, err.Error())
			return
		}
		utils.FailedResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Nickname updated", res)
}

// GetPendingFriendRequest godoc
//
//	@Summary	Get pending & declined friend request by user wallet address
//
// Description get pending & declined friend request by user wallet address
//
//	@Tags		Friend
//	@Accept		json
//	@Produce	json
//	@Param		user_wallet_address	path string true "user wallet addres"
//	@Success	200		{object}	dtos.PendingFriendResponse
//	@Failure	400		"Invalid Request"
//	@Failure	500		"Internal Server Error"
//	@Router		/friends/get-pending-request/{user_wallet_address} [get]
func GetPendingFriendRequest(c *gin.Context) {
	userWalletAddress := c.Param("user_wallet_address")

	res, err := services.GetPendingFriendRequestService(userWalletAddress)
	if err != nil {
		if err.Error() == "invalid request" {
			utils.FailedResponse(c, http.StatusBadRequest, err.Error())
		} else {
			utils.FailedResponse(c, http.StatusInternalServerError, "Failed to fetch pending friend request")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success fetch pending friend request", res)
}

// GetPendingFriendRequestByFriendWalletAddress godoc
//
//	@Summary	Get pending & declined friend request by friend wallet address
//
// Description get pending & declined friend request by friend wallet address
//
//	@Tags		Friend
//	@Accept		json
//	@Produce	json
//	@Param		friend_wallet_address	path string true "friend wallet addres"
//	@Success	200		{object}	dtos.PendingFriendResponse
//	@Failure	400		"Invalid Request"
//	@Failure	500		"Internal Server Error"
//	@Router		/friends/get-pending-request-by-friend/{friend_wallet_address} [get]
func GetPendingFriendRequestByFriendWalletAddress(c *gin.Context) {
	friendWalletAddress := c.Param("friend_wallet_address")

	res, err := services.GetPendingFriendRequestServiceRequestedUser(friendWalletAddress)
	if err != nil {
		if err.Error() == "invalid request" {
			utils.FailedResponse(c, http.StatusBadRequest, err.Error())
		} else {
			utils.FailedResponse(c, http.StatusInternalServerError, "Failed to fetch pending friend request")
		}
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Success fetch pending friend request", res)
}
