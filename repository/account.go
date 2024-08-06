package repository

import (
	"fmt"
	"welloff-bank/model"

	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	Pg *sqlx.DB
}

func (ac *AccountRepository) CreateAccount(user_id string, name string, status string) error {
	_, err := ac.Pg.Exec(
		`INSERT INTO "account" (user_id, name, status)
		VALUES ($1, $2, $3)
		`,
		user_id,
		name,
		status,
	)

	return err
}

func (ac *AccountRepository) GetAccount(acc_id string) (*model.Account, error) {
	account := new(model.Account)
	err := ac.Pg.Get(
		account,
		`SELECT acc.id, acc.user_id, acc.name, acc.status, acc.created_at, acc.updated_at 
		FROM "account" acc WHERE acc.id = $1`,
		acc_id,
	)

	return account, err
}

func (ac *AccountRepository) GetMyAccounts(user_id string, limit int, offset int) (*[]model.Account, error) {
	accounts := new([]model.Account)
	err := ac.Pg.Select(
		accounts,
		`
		SELECT 
			acc.id, acc.user_id, acc.name, acc.status, acc.created_at, acc.updated_at 
		FROM 
			"account" acc 
		WHERE 
			acc.user_id = $1
		ORDER BY
			acc.status
		LIMIT 
			$2
		OFFSET 
			$3
		`,
		user_id,
		limit,
		offset,
	)

	return accounts, err
}

func (ac *AccountRepository) DisableAccount(acc_id string) error {
	_, err := ac.Pg.Exec(
		`UPDATE "account"
		SET status = 'inactive'
		WHERE id = $1`,
		acc_id,
	)

	return err
}

func (ac *AccountRepository) GetBalanceSnapshot(account_id string) (*model.AccountBalance, error) {
	balance_snapshot := new(model.AccountBalance)
	err := ac.Pg.Get(
		balance_snapshot,
		`SELECT bs.account_id, bs.balance, bs.updated_at FROM "balance_snapshot" bs WHERE bs.account_id = $1`,
		account_id,
	)

	return balance_snapshot, err
}

func (ac *AccountRepository) GetAccounts(limit int, offset int) (*[]model.Account, error) {
	accounts := new([]model.Account)
	err := ac.Pg.Select(
		accounts,
		`
		SELECT 
			acc.id, acc.user_id, acc.name, acc.status, acc.created_at, acc.updated_at 
		FROM 
			"account" acc 
		ORDER BY
			acc.created_at DESC
		LIMIT 
			$1
		OFFSET 
			$2
		`,
		limit,
		offset,
	)

	return accounts, err
}

func (ac *AccountRepository) BulkUpsertBalanceSnapshots(balances *[]model.AccountBalance) error {
	query := `INSERT INTO "balance_snapshot" (account_id, balance, created_at, updated_at) VALUES `
	values := []interface{}{}

	for i, balance := range *balances {
		num := i * 4
		query += fmt.Sprintf("($%d, $%d, $%d, $%d),", num+1, num+2, num+3, num+4)
		values = append(values, balance.AccountId, balance.Balance, balance.Date, balance.Date)
	}

	query = query[:len(query)-1]
	query += " ON CONFLICT (account_id) DO UPDATE SET balance = EXCLUDED.balance, updated_at = EXCLUDED.updated_at"

	_, err := ac.Pg.Exec(query, values...)

	return err
}
