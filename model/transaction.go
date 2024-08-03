package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Transaction struct {
	Id uuid.UUID `db:"id" json:"id"`
	// 'deposit' | 'withdrawal' | 'transfer' | 'refund'
	Kind          string          `db:"kind" json:"kind"`
	FromAccountId *uuid.UUID      `db:"from_account_id" json:"from_account_id"`
	ToAccountId   *uuid.UUID      `db:"to_account_id" json:"to_account_id"`
	Amount        decimal.Decimal `db:"amount" json:"amount"`
	DateIssued    time.Time       `db:"date_issued" json:"date_issued"`
	RelatedTransactionId *uuid.UUID `db:"related_transaction_id" json:"related_transaction_id"`
}