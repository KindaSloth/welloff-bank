package utils

import (
	"encoding/json"
	"welloff-bank/model"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func GetUser(ctx *gin.Context) (*model.User, error) {
	userJson := ctx.GetString("user")
	var user model.User
	err := json.Unmarshal([]byte(userJson), &user)

	return &user, err
}

func GetAccountBalance(ctx *gin.Context, account_id uuid.UUID) (model.AccountBalance, error) {
	return model.AccountBalance{
		AccountId: account_id,
		Balance:   decimal.NewFromInt(0),
	}, nil
}

func FilterTransactionsByKind(transactions []model.Transaction, kind string) []model.Transaction {
	var filtered_transactions []model.Transaction

	for _, transaction := range transactions {
		if transaction.Kind == kind {
			filtered_transactions = append(filtered_transactions, transaction)
		}
	}

	return filtered_transactions
}
