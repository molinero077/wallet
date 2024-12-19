package model

type Operation struct {
	WalletId      string  `json:"valletId"`
	OperationType string  `json:"operationType"`
	Amount        float32 `json:"amount"`
}

type WalletBalance struct {
	WalletId string  `json:"valletId"`
	Amount   float32 `json:"amount"`
}

type Wallet struct {
	WalletId     string `json:"valletId"`
	UserId       string `json:"userId"`
	CreationDate string `json:"creationDate"`
}

type WalletProvider interface {
	GetBalance(id string) (float32, error)
	CarryOperation(ops Operation) error
}
