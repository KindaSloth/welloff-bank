package repository

import (
	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	Pg *sqlx.DB
}