package dto

// help me generate user request the with body user id and wallet balance
type UserDepositRequest struct {
	UserID        string  `json:"user_id"`
	DepositAmount float64 `json:"deposit_amount"`
}

type UserDepositResponse struct {
	UserID       string  `json:"user_id"`
	WalletAmount float64 `json:"wallet_amount"`
}
