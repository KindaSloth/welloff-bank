package repository

import (
	"welloff-bank/model"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	Pg *sqlx.DB
}

func (ur *UserRepository) CreateUser(email string, password string) error {
	_, err := ur.Pg.Exec(
		`INSERT INTO "user" (email, password)
		VALUES ($1, $2)`,
		email,
		password,
	)

	return err
}

func (ur *UserRepository) GetUserById(id uuid.UUID) (*model.User, error) {
	user := new(model.User)
	err := ur.Pg.Get(
		user,
		`SELECT u.id, u.email, u.password, u.created_at, u.updated_at
		FROM "user" u WHERE u.id=$1`,
		id,
	)

	return user, err
}

func (ur *UserRepository) GetUserByEmail(email string) (*model.User, error) {
	user := new(model.User)
	err := ur.Pg.Get(
		user,
		`SELECT u.id, u.email, u.password, u.created_at, u.updated_at
		FROM "user" u WHERE u.email=$1`,
		email,
	)

	return user, err
}
