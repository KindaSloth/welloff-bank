package repository

import (
	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	Pg *sqlx.DB
}