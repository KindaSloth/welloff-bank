package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id                uuid.UUID `db:"id" json:"id"`
	Email             string    `db:"email" json:"email"`
	EncryptedPassword string    `db:"password" json:"-"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}