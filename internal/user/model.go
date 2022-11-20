package user

type User struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type TransactionRequest struct {
	UserID          string `json:"user_id"`
	Amount          int    `json:"amount"`
	TransactionType string `json:"transaction_type"`
}

type UserTransaction struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Type      string `json:"type"`
	Status    string `json:"status"`
	Amount    int    `json:"amount"`
	CreatedAt string `json:"created_at"`
}
