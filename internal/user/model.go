package user

type User struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Balance int    `json:"balance"`
}

type UserRequest struct {
	Id            string  `json:"id"`
	Amount        float64 `json:"amount"`
	OperationType string  `json:"operation_type"`
}
