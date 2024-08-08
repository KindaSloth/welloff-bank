package server

import (
	"context"
	"log"
	"strconv"
	"welloff-bank/utils"

	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
)

type CreateAccountRequest struct {
	Name string `json:"name"`
}

func (s *Server) CreateAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := CreateAccountRequest{}
		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [CreateAccount] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		err = s.Repositories.AccountRepository.CreateAccount(user.Id.String(), req.Name, "active")
		if err != nil {
			log.Println("[ERROR] [CreateAccount] failed to create account: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to create account"})
			return
		}

		ctx.Status(200)
	}
}

type GetAccountResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Balance string `json:"balance"`
	// 'active' | 'inactive'
	Status string `json:"status"`
}

func (s *Server) GetAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account_id := ctx.Param("id")
		if account_id == "" {
			ctx.JSON(400, gin.H{"error": "Missing account id"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [GetAccount] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		account, err := s.Repositories.AccountRepository.GetAccount(account_id)
		if err != nil {
			ctx.JSON(404, gin.H{"error": "Account not found"})
			return
		}

		if account.UserId.String() != user.Id.String() {
			ctx.JSON(401, gin.H{"error": "User is not the owner of the account"})
			return
		}

		account_balance, err := utils.GetAccountBalance(context.Background(), account.Id, s.Repositories, true)
		if err != nil {
			log.Println("[ERROR] [GetAccount] failed to get account balance: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to get account balance"})
			return
		}

		ctx.JSON(200, gin.H{"payload": GetAccountResponse{
			Id:      account.Id.String(),
			Name:    account.Name,
			Balance: account_balance.Balance.String(),
			Status:  account.Status,
		}})
	}
}

type GetAccountsResponse struct {
	AccountId string `json:"account_id"`
	Name      string `json:"name"`
	// 'active' | 'inactive'
	Status string `json:"status"`
}

func (s *Server) GetAccounts() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limit, err := strconv.Atoi(ctx.Query("limit"))
		if err != nil || limit == 0 {
			limit = 10
		}

		offset, err := strconv.Atoi(ctx.Query("offset"))
		if err != nil || offset < 0 {
			offset = 0
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [GetAccount] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		accounts, err := s.Repositories.AccountRepository.GetMyAccounts(user.Id.String(), limit, offset)
		if err != nil {
			log.Println("[ERROR] [GetAccounts] failed to get accounts: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to get accounts"})
			return
		}

		var payload []GetAccountsResponse

		for _, account := range *accounts {
			payload = append(payload, GetAccountsResponse{
				AccountId: account.Id.String(),
				Name:      account.Name,
				Status:    account.Status,
			})
		}

		ctx.JSON(200, gin.H{"payload": payload})
	}
}

func (s *Server) DisableAccount() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		account_id := ctx.Param("id")
		if account_id == "" {
			ctx.JSON(400, gin.H{"error": "Missing id param"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [DisableAccount] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		account, err := s.Repositories.AccountRepository.GetAccount(account_id)
		if err != nil {
			log.Println("[ERROR] [DisableAccount] failed to get account: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to get account"})
			return
		}

		if account.UserId != user.Id {
			ctx.JSON(500, gin.H{"error": "This account does not belongs to the user"})
			return
		}

		account_balance, err := utils.GetAccountBalance(context.Background(), account.Id, s.Repositories, false)
		if err != nil {
			log.Println("[ERROR] [DisableAccount] failed to get account balance: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to get account balance"})
			return
		}

		if account_balance.Balance.GreaterThan(decimal.NewFromInt(0)) {
			ctx.JSON(500, gin.H{"error": "Account still has balance and cannot be deleted"})
			return
		}

		err = s.Repositories.AccountRepository.DisableAccount(account_id)
		if err != nil {
			log.Println("[ERROR] [DisableAccount] failed to disable account: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to disable account"})
			return
		}

		ctx.Status(200)
	}
}
