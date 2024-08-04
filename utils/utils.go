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
