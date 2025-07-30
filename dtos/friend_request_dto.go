package dtos

type AddFriendRequest struct {
	UserWalletAddress   string `json:"user_wallet_address" binding:"required"`
	FriendWalletAddress string `json:"friend_wallet_address" binding:"required"`
}

type AddFriendResponse struct {
	ID                  string `json:"id"`
	UserWalletAddress   string `json:"user_wallet_address" binding:"required"`
	FriendWalletAddress string `json:"friend_wallet_address" binding:"required"`
	Status              string `json:"status"`
}

type AcceptFriendRequest struct {
	ID string `json:"id" binding:"required"`
}

type AcceptFriendResponse struct {
	ID                  string  `json:"id"`
	UserWalletAddress   string  `json:"user_wallet_address"`
	FriendWalletAddress string  `json:"friend_wallet_address"`
	Nickname            *string `json:"nickname"`
}

type DeclineFriendRequest struct {
	ID string `json:"id" binding:"required"`
}

type DeclineFriendResponse struct {
	ID                  string `json:"id"`
	UserWalletAddress   string `json:"user_wallet_address" binding:"required"`
	FriendWalletAddress string `json:"friend_wallet_address" binding:"required"`
	Status              string `json:"status"`
}

type PendingFriendResponse struct {
	ID                  string `json:"id"`
	UserWalletAddress   string `json:"user_wallet_address"`
	FriendWalletAddress string `json:"friend_wallet_address"`
	Status              string `json:"status"`
}
