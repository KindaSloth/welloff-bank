package repository

import (
	"fmt"
	"log"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/valkey-io/valkey-go"

	_ "github.com/lib/pq"
)

type Repositories struct {
	Pg                    *sqlx.DB
	Valkey                valkey.Client
	UserRepository        UserRepository
	AccountRepository     AccountRepository
	TransactionRepository TransactionRepository
}

func New() Repositories {
	pg_user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		log.Fatal("Missing POSTGRES_USER env")
	}
	pg_dbname, ok := os.LookupEnv("POSTGRES_DBNAME")
	if !ok {
		log.Fatal("Missing POSTGRES_DBNAME env")
	}
	pg_password, ok := os.LookupEnv("POSTGRES_PASSWORD")
	if !ok {
		log.Fatal("Missing POSTGRES_PASSWORD env")
	}
	pg_sslmode, ok := os.LookupEnv("POSTGRES_SSLMODE")
	if !ok {
		log.Fatal("Missing POSTGRES_SSLMODE env")
	}
	valkey_addr, ok := os.LookupEnv("VALKEY_ADDRESS")
	if !ok {
		log.Fatal("Missing VALKEY_ADDRESS env")
	}

	pg, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s password=%s sslmode=%s", pg_user, pg_dbname, pg_password, pg_sslmode))
	if err != nil {
		msg := fmt.Sprintf("[ERROR] failed to create database: %s", err)
		log.Fatal(msg)
	}

	log.Println("Connected to Postgres")

	pg.SetMaxOpenConns(200)

	valkey, err := valkey.NewClient(valkey.ClientOption{InitAddress: []string{valkey_addr}})
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to Valkey")

	return Repositories{
		Pg:                    pg,
		Valkey:                valkey,
		UserRepository:        UserRepository{pg},
		AccountRepository:     AccountRepository{pg},
		TransactionRepository: TransactionRepository{pg},
	}
}
