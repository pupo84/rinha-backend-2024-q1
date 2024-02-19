package dto

type Transaction struct {
	Amount      int64  `json:"valor"`
	Nature      string `json:"tipo"`
	Description string `json:"descricao"`
	CreatedAt   string `json:"realizada_em"`
}

type Balance struct {
	Limit     int64  `json:"limite"`
	Balance   int64  `json:"saldo"`
	UpdatedAt string `json:",omitempty"`
}

type TransactionErrorResponse struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type CreateUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Document string `json:"document"`
	Limit    int64  `json:"limit"`
}

type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Document string `json:"document"`
	Limit    int64  `json:"limit"`
}
type BalanceResponse struct {
	Total         int64  `json:"total"`
	StatementDate string `json:"data_extrato"`
	Limit         int64  `json:"limite"`
}

type TransactionResponse struct {
	BalanceResponse `json:"saldo"`
	Transactions    []Transaction `json:"ultimas_transacoes"`
}
