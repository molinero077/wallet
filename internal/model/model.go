package model

import (
	"errors"
	"math"
	"strings"
)

var (
	ErrWalletIdMissing     error = errors.New("walletId is required")
	ErrOperationTypeMissin error = errors.New("operationType is required")
	ErrAmountMissing       error = errors.New("amount is required")
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

func (op *Operation) CheckRequiredFields() error {
	if op.WalletId == "" {
		return ErrWalletIdMissing
	}
	if op.OperationType == "" {
		return ErrOperationTypeMissin
	}
	if op.Amount == 0 {
		return ErrAmountMissing
	}

	return nil
}

type WalletBalance struct {
	WalletId string  `json:"valletId"`
	Balance  float64 `json:"balance"`
}

type WalletProvider interface {
	GetWalletBalance(id string) (*WalletBalance, error)
	WalletOperation(ops Operation) error
}

var (
	ErrNonExistentWallet error = errors.New("non-existent wallet")
	ErrZeroAmount        error = errors.New("zero amount")
)
