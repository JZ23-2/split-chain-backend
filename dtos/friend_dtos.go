package dtos

type GetFriendResponse struct {
	ID                  string  `json:"id"`
	FriendWalletAddress string  `json:"friend_wallet_address"`
	Nickname            *string `json:"nickname"`
}
