package repository

import (
	"time"
	"welloff-bank/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type TransactionRepository struct {
	Pg *sqlx.DB
}

func (tr *TransactionRepository) GetTransaction(transaction_id string) (*model.Transaction, error) {
	transaction := new(model.Transaction)
	err := tr.Pg.Get(
		transaction,
		`SELECT tx.id, tx.kind, tx.from_account_id, tx.to_account_id, tx.amount, tx.date_issued, tx.related_transaction_id
		FROM "transaction" tx WHERE tx.id = $1`,
		transaction_id,
	)

	return transaction, err
}

func (tr *TransactionRepository) GetTransactionsByAccount(account_id string, limit int, offset int) (*[]model.Transaction, error) {
	transactions := new([]model.Transaction)
	err := tr.Pg.Select(
		transactions,
		`
		SELECT 
			tx.id, tx.kind, tx.from_account_id, tx.to_account_id, tx.amount, tx.date_issued, tx.related_transaction_id
		FROM 
			"transaction" tx 
		WHERE 
			tx.from_account_id = $1
		ORDER BY
			tx.date_issued DESC
		LIMIT 
			$2
		OFFSET 
			$3
		`,
		account_id,
		limit,
		offset,
	)

	return transactions, err
}

func (tr *TransactionRepository) GetAllTransactionsByAccount(account_id string) (*[]model.Transaction, error) {
	transactions := new([]model.Transaction)
	err := tr.Pg.Select(
		transactions,
		`
		SELECT 
			tx.id, tx.kind, tx.from_account_id, tx.to_account_id, tx.amount, tx.date_issued, tx.related_transaction_id
		FROM 
			"transaction" tx 
		WHERE 
			(tx.from_account_id = $1 OR tx.to_account_id = $1)
		ORDER BY
			tx.date_issued DESC
		`,
		account_id,
	)

	return transactions, err
}

func (tr *TransactionRepository) GetTransactionsByDate(account_id string, date_from time.Time, date_to time.Time) (*[]model.Transaction, error) {
	transactions := new([]model.Transaction)
	err := tr.Pg.Select(
		transactions,
		`
		SELECT 
			tx.id, tx.kind, tx.from_account_id, tx.to_account_id, tx.amount, tx.date_issued, tx.related_transaction_id
		FROM 
			"transaction" tx 
		WHERE 
			(tx.from_account_id = $1 OR tx.to_account_id = $1)
		AND 
			tx.date_issued BETWEEN $2 AND $3
		ORDER BY
			tx.date_issued DESC
		`,
		account_id,
		date_from,
		date_to,
	)

	return transactions, err
}

func (tr *TransactionRepository) CreateTransaction(transaction_id uuid.UUID, kind string, from_account_id *string, to_account_id *string, amount decimal.Decimal, related_transaction_id *string) error {
	_, err := tr.Pg.Exec(
		`INSERT INTO "transaction" (id, kind, from_account_id, to_account_id, amount, related_transaction_id)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		transaction_id,
		kind,
		from_account_id,
		to_account_id,
		amount,
		related_transaction_id,
	)

	return err
}
