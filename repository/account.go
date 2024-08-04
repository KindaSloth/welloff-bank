package repository

import (
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
