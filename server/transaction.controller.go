package server

import (
	"context"
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

		account_balance, err := utils.GetAccountBalance(context.Background(), account.Id, s.Repositories, false)
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

		account_balance, err := utils.GetAccountBalance(context.Background(), account.Id, s.Repositories, false)
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

		transaction, err := s.Repositories.TransactionRepository.GetTransaction(transaction_id)
		if err != nil {
			log.Println("[ERROR] [RefundTransaction] failed to get transaction: ", err)
			ctx.JSON(404, gin.H{"error": "Transaction not found"})
			return
		}

		if transaction.Kind != "transfer" {
			ctx.JSON(400, gin.H{"error": "Transaction cannot be refunded"})
			return
		}

		user, err := utils.GetUser(ctx)
		if err != nil {
			log.Println("[ERROR] [RefundTransaction] failed to get user from context: ", err)
			ctx.Status(401)
			return
		}

		from_account, err := s.Repositories.AccountRepository.GetAccount(transaction.FromAccountId.String())
		if err != nil {
			log.Println("[ERROR] [RefundTransaction] failed to get account: ", err)
			ctx.JSON(404, gin.H{"error": "Account not found"})
			return
		}

		to_account, err := s.Repositories.AccountRepository.GetAccount(transaction.ToAccountId.String())
		if err != nil {
			log.Println("[ERROR] [RefundTransaction] failed to get account: ", err)
			ctx.JSON(404, gin.H{"error": "Account not found"})
			return
		}

		if (from_account.UserId.String() != user.Id.String()) && (to_account.UserId.String() != user.Id.String()) {
			ctx.JSON(401, gin.H{"error": "User is not the owner of any of the accounts"})
			return
		}

		to_account_balance, err := utils.GetAccountBalance(context.Background(), to_account.Id, s.Repositories, false)
		if err != nil {
			log.Println("[ERROR] [RefundTransaction] failed to get account balance: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete refund transaction"})
			return
		}

		if to_account_balance.Balance.LessThan(transaction.Amount) {
			ctx.JSON(400, gin.H{"error": "Insufficient balance"})
			return
		}

		refund_transaction_id, err := uuid.NewV7()
		if err != nil {
			log.Println("[ERROR] [TransferTransaction] failed to create transaction id: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete transfer transaction"})
			return
		}

		tx, err := s.Repositories.TransactionRepository.GetTransaction(refund_transaction_id.String())
		if tx != nil && err == nil {
			ctx.JSON(500, gin.H{"error": "Duplicated transaction"})
			return
		}

		from_account_id := transaction.FromAccountId.String()
		to_account_id := transaction.ToAccountId.String()

		err = s.Repositories.TransactionRepository.CreateTransaction(refund_transaction_id, "refund", &from_account_id, &to_account_id, transaction.Amount, nil)
		if err != nil {
			log.Println("[ERROR] [RefundTransaction] failed to create transaction: ", err)
			ctx.JSON(500, gin.H{"error": "Failed to complete refund transaction"})
			return
		}

		ctx.Status(200)
	}
}
