package model

import (
	"errors"
	"math"
	"strings"
)

type Operation struct {
	WalletId      string  `json:"valletId"`
	OperationType string  `json:"operationType"`
	Amount        float64 `json:"amount"`
}

// GetAmount - возращает сумму операции. DEPOSIT +сумма, WITHDRAW -сумма
func (op *Operation) GetAmount() float64 {
	switch strings.ToLower(op.OperationType) {
	case "deposit":
		return math.Abs(op.Amount)
	case "withdraw":
		return math.Abs(op.Amount) * -1
	default:
		return 0
	}
}

type WalletBalance struct {
	WalletId string  `json:"valletId"`
	Balance  float64 `json:"balance"`
}

type WalletProvider interface {
	GetBalance(id string) (*WalletBalance, error)
	CarryOperation(ops Operation) error
}

var (
	ErrNonExistentWallet error = errors.New("non-existent wallet")
	ErrZeroAmount        error = errors.New("zero amount")
)
