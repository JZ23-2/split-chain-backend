package services

import (
	"errors"

	"github.com/JZ23-2/splitbill-backend/database"
	"github.com/JZ23-2/splitbill-backend/dtos"
	"github.com/JZ23-2/splitbill-backend/models"
	"gorm.io/gorm"
)

func AcceptFriendRequestService(req dtos.AcceptFriendRequest) (dtos.AcceptFriendResponse, dtos.AcceptFriendResponse, error) {
	var pending models.PendingFriendRequest

	if err := database.DB.First(&pending, "id = ?", req.ID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return dtos.AcceptFriendResponse{}, dtos.AcceptFriendResponse{}, err
		}
		return dtos.AcceptFriendResponse{}, dtos.AcceptFriendResponse{}, err
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
		return dtos.AcceptFriendResponse{}, dtos.AcceptFriendResponse{}, err
	}

	if err := database.DB.Create(&friend2).Error; err != nil {
		return dtos.AcceptFriendResponse{}, dtos.AcceptFriendResponse{}, err
	}

	_ = database.DB.Delete(&pending)

	resp1 := dtos.AcceptFriendResponse{
		ID:                  friend1.ID,
		UserWalletAddress:   friend1.UserWalletAddress,
		FriendWalletAddress: friend1.FriendWalletAddress,
		Nickname:            &friend1.Nickname,
	}

	resp2 := dtos.AcceptFriendResponse{
		ID:                  friend2.ID,
		UserWalletAddress:   friend2.UserWalletAddress,
		FriendWalletAddress: friend2.FriendWalletAddress,
		Nickname:            &friend2.Nickname,
	}

	return resp1, resp2, nil
}

func DeclineFriendRequestService(req dtos.DeclineFriendRequest) (dtos.DeclineFriendResponse, error) {
	var pending models.PendingFriendRequest

	if err := database.DB.First(&pending, "id = ?", req.ID).Error; err != nil {
		return dtos.DeclineFriendResponse{}, err
	}

	pending.Status = "Declined"
	if err := database.DB.Save(&pending).Error; err != nil {
		return dtos.DeclineFriendResponse{}, err
	}

	response := dtos.DeclineFriendResponse{
		ID:                  pending.ID,
		UserWalletAddress:   pending.UserWalletAddress,
		FriendWalletAddress: pending.FriendWalletAddress,
		Status:              pending.Status,
	}

	return response, nil

}

func AddFriendRequestService(req dtos.AddFriendRequest) (any, string, int, error) {
	var user, friend models.User

	if err := database.DB.First(&user, "wallet_address = ?", req.UserWalletAddress).Error; err != nil {
		return nil, "User Not Found", 404, err
	}

	if err := database.DB.First(&friend, "wallet_address = ?", req.FriendWalletAddress).Error; err != nil {
		return nil, "Friend Not Found", 404, err
	}

	var existing models.PendingFriendRequest
	if err := database.DB.
		Where("user_wallet_address = ? AND friend_wallet_address = ? AND status = ?", req.UserWalletAddress, req.FriendWalletAddress, "Declined").
		First(&existing).Error; err == nil {

		existing.Status = "Pending"
		if err := database.DB.Save(&existing).Error; err != nil {
			return nil, "Failed to re-activate declined request", 500, err
		}
		return existing, "Successfully Added Friend Request", 200, nil
	}

	if err := database.DB.
		Where(
			"(user_wallet_address = ? AND friend_wallet_address = ?) OR (user_wallet_address = ? AND friend_wallet_address = ?)",
			req.UserWalletAddress, req.FriendWalletAddress,
			req.FriendWalletAddress, req.UserWalletAddress,
		).
		First(&existing).Error; err == nil {

		// Accept the request and make it a Friend record
		friend1 := models.Friend{
			UserWalletAddress:   req.FriendWalletAddress,
			FriendWalletAddress: req.UserWalletAddress,
		}
		friend2 := models.Friend{
			UserWalletAddress:   req.UserWalletAddress,
			FriendWalletAddress: req.FriendWalletAddress,
		}

		if err := database.DB.Create(&friend1).Error; err != nil {
			return nil, "Failed to create friends", 500, err
		}
		if err := database.DB.Create(&friend2).Error; err != nil {
			return nil, "Failed to create friends", 500, err
		}

		_ = database.DB.
			Where(
				"(user_wallet_address = ? AND friend_wallet_address = ?) OR (user_wallet_address = ? AND friend_wallet_address = ?)",
				req.UserWalletAddress, req.FriendWalletAddress,
				req.FriendWalletAddress, req.UserWalletAddress,
			).
			Delete(&models.PendingFriendRequest{})

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

		return map[string]any{
			"friend_1": friend1Resp,
			"friend_2": friend2Resp,
		}, "Friend Request Accepted", 200, nil
	}

	var friendCheck models.Friend
	if err := database.DB.
		Where("user_wallet_address = ? AND friend_wallet_address = ?", req.UserWalletAddress, req.FriendWalletAddress).
		First(&friendCheck).Error; err == nil {
		return nil, "Already Friend", 409, errors.New("already friends")
	}

	newRequest := models.PendingFriendRequest{
		UserWalletAddress:   req.UserWalletAddress,
		FriendWalletAddress: req.FriendWalletAddress,
	}
	if err := database.DB.Create(&newRequest).Error; err != nil {
		return nil, "Failed to create friend request", 500, err
	}

	response := dtos.AddFriendResponse{
		ID:                  newRequest.ID,
		UserWalletAddress:   newRequest.UserWalletAddress,
		FriendWalletAddress: newRequest.FriendWalletAddress,
		Status:              "Pending",
	}

	return response, "Successfully Added Friend Request", 200, nil
}

func GetFriendService(userWalletAddress string) ([]dtos.FriendResponse, error) {
	var friends []dtos.FriendResponse

	err := database.DB.
		Table("friends").
		Select("id, nickname, friend_wallet_address").
		Where("user_wallet_address = ?", userWalletAddress).
		Find(&friends).Error

	return friends, err
}

func AddFriendNicknameService(req dtos.FriendNicknameRequest) (*dtos.FriendResponse, error) {
	var friend models.Friend

	err := database.DB.
		Where("user_wallet_address = ? AND friend_wallet_address = ?", req.UserWalletAddress, req.FriendWalletAddress).
		First(&friend).Error

	if err != nil {
		return nil, errors.New("user or friend not found")
	}

	friend.Nickname = *req.Nickname
	if err := database.DB.Save(&friend).Error; err != nil {
		return nil, errors.New("failed to update nickname")
	}

	res := &dtos.FriendResponse{
		ID:                  friend.ID,
		FriendWalletAddress: friend.FriendWalletAddress,
		Nickname:            &friend.Nickname,
	}

	return res, nil
}

func GetPendingFriendRequestService(userWalletAddress string) ([]dtos.PendingFriendResponse, error) {
	if userWalletAddress == "" {
		return nil, errors.New("invalid request")
	}

	var pendingRequests []models.PendingFriendRequest
	err := database.DB.
		Where("user_wallet_address = ? AND status = ?", userWalletAddress, "Pending").
		Find(&pendingRequests).Error

	if err != nil {
		return nil, err
	}

	if len(pendingRequests) == 0 {
		return []dtos.PendingFriendResponse{}, nil
	} else {
		var responses []dtos.PendingFriendResponse
		for _, r := range pendingRequests {
			responses = append(responses, dtos.PendingFriendResponse{
				ID:                  r.ID,
				UserWalletAddress:   r.UserWalletAddress,
				FriendWalletAddress: r.FriendWalletAddress,
				Status:              r.Status,
			})
		}

		return responses, nil

	}
}
