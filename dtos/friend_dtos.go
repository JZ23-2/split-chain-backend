package dtos

type FriendResponse struct {
	ID                  string  `json:"id"`
	FriendWalletAddress string  `json:"friend_wallet_address"`
	Nickname            *string `json:"nickname"`
}

type FriendNicknameRequest struct {
	UserWalletAddress   string  `json:"user_wallet_address" binding:"required"`
	FriendWalletAddress string  `json:"friend_wallet_address" binding:"required"`
	Nickname            *string `json:"nickname" binding:"required"`
}
