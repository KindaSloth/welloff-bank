package repository

import (
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	Pg *sqlx.DB
}