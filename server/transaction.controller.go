package server

import (
	"log"
	"welloff-bank/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func (s *Server) GetTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		if id == "" {
			ctx.JSON(400, gin.H{"error": "Missing id param"})
			return
		}

		transaction, err := s.Repositories.TransactionRepository.GetTransaction(id)
		if err != nil {
			ctx.JSON(404, gin.H{"error": "Transaction not found"})
			return
		}

		ctx.JSON(200, gin.H{"payload": transaction})
	}
}

type DepositTransactionRequest struct {
	Amount      decimal.Decimal `json:"amount"`
	ToAccountId string          `json:"to_account_id"`
}

func (s *Server) DepositTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := DepositTransactionRequest{}
		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [DepositTransaction] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		account, err := s.Repositories.AccountRepository.GetAccount(req.ToAccountId)
		if err != nil {
			log.Println("[ERROR] [DepositTransaction] failed to get account: ", err)
			ctx.JSON(404, gin.H{"error": "Account not found"})
			return
		}

		if account.UserId.String() != user.Id.String() {
			ctx.JSON(401, gin.H{"error": "User is not the owner of the account"})
			return
		}

		transaction_id, err := uuid.NewV7()
		if err != nil {
			log.Println("[ERROR] [DepositTransaction] failed to create transaction id: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete deposit transaction"})
			return
		}

		tx, err := s.Repositories.TransactionRepository.GetTransaction(transaction_id.String())
		if tx != nil && err == nil {
			ctx.JSON(500, gin.H{"error": "Duplicated transaction"})
			return
		}

		err = s.Repositories.TransactionRepository.CreateTransaction(transaction_id, "deposit", nil, &req.ToAccountId, req.Amount, nil)
		if err != nil {
			log.Println("[ERROR] [DepositTransaction] failed to create transaction: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete deposit transaction"})
			return
		}

		ctx.Status(200)
	}
}

type WithdrawalTransactionRequest struct {
	Amount        decimal.Decimal `json:"amount"`
	FromAccountId string          `json:"from_account_id"`
}

func (s *Server) WithdrawalTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := WithdrawalTransactionRequest{}
		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [WithdrawalTransaction] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		account, err := s.Repositories.AccountRepository.GetAccount(req.FromAccountId)
		if err != nil {
			log.Println("[ERROR] [WithdrawalTransaction] failed to get account: ", err)
			ctx.JSON(404, gin.H{"error": "Account not found"})
			return
		}

		if account.UserId.String() != user.Id.String() {
			ctx.JSON(401, gin.H{"error": "User is not the owner of the account"})
			return
		}

		account_balance, err := utils.GetAccountBalance(ctx, account.Id, s.Repositories)
		if err != nil {
			log.Println("[ERROR] [WithdrawalTransaction] failed to get account balance: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete withdrawal transaction"})
			return
		}

		if account_balance.Balance.LessThan(req.Amount) {
			ctx.JSON(400, gin.H{"error": "Insufficient balance"})
			return
		}

		transaction_id, err := uuid.NewV7()
		if err != nil {
			log.Println("[ERROR] [WithdrawalTransaction] failed to create transaction id: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete withdrawal transaction"})
			return
		}

		tx, err := s.Repositories.TransactionRepository.GetTransaction(transaction_id.String())
		if tx != nil && err == nil {
			ctx.JSON(500, gin.H{"error": "Duplicated transaction"})
			return
		}

		err = s.Repositories.TransactionRepository.CreateTransaction(transaction_id, "withdrawal", &req.FromAccountId, nil, req.Amount, nil)
		if err != nil {
			log.Println("[ERROR] [WithdrawalTransaction] failed to create transaction: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete withdrawal transaction"})
			return
		}

		ctx.Status(200)
	}
}

type TransferTransactionRequest struct {
	Amount        decimal.Decimal `json:"amount"`
	FromAccountId string          `json:"from_account_id"`
	ToAccountId   string          `json:"to_account_id"`
}

func (s *Server) TransferTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		req := TransferTransactionRequest{}
		if ctx.ShouldBindJSON(&req) != nil {
			ctx.JSON(422, gin.H{"error": "Invalid input"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [TransferTransaction] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		account, err := s.Repositories.AccountRepository.GetAccount(req.FromAccountId)
		if err != nil {
			log.Println("[ERROR] [TransferTransaction] failed to get account: ", err)
			ctx.JSON(404, gin.H{"error": "Account not found"})
			return
		}

		if account.UserId.String() != user.Id.String() {
			ctx.JSON(401, gin.H{"error": "User is not the owner of the account"})
			return
		}

		account_balance, err := utils.GetAccountBalance(ctx, account.Id, s.Repositories)
		if err != nil {
			log.Println("[ERROR] [TransferTransaction] failed to get account balance: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete transfer transaction"})
			return
		}

		if account_balance.Balance.LessThan(req.Amount) {
			ctx.JSON(400, gin.H{"error": "Insufficient balance"})
			return
		}

		transaction_id, err := uuid.NewV7()
		if err != nil {
			log.Println("[ERROR] [TransferTransaction] failed to create transaction id: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete transfer transaction"})
			return
		}

		tx, err := s.Repositories.TransactionRepository.GetTransaction(transaction_id.String())
		if tx != nil && err == nil {
			ctx.JSON(500, gin.H{"error": "Duplicated transaction"})
			return
		}

		err = s.Repositories.TransactionRepository.CreateTransaction(transaction_id, "transfer", &req.FromAccountId, &req.ToAccountId, req.Amount, nil)
		if err != nil {
			log.Println("[ERROR] [TransferTransaction] failed to create transaction: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete transfer transaction"})
			return
		}

		ctx.Status(200)
	}
}

func (s *Server) RefundTransaction() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		transaction_id := ctx.Param("id")
		if transaction_id == "" {
			ctx.JSON(400, gin.H{"error": "Missing id param"})
			return
		}

		// TODO
		ctx.Status(400)
	}
}
