package utils

import (
	"context"
	"encoding/json"
	"errors"
	"time"
	"welloff-bank/model"
	"welloff-bank/repository"

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

func GetAccountBalance(ctx context.Context, account_id uuid.UUID, repostiories repository.Repositories) (*model.AccountBalance, error) {
	account, err := repostiories.AccountRepository.GetAccount(account_id.String())
	if err != nil {
		return nil, errors.New("failed to get account")
	}

	balance := decimal.Zero
	now := time.Now().UTC()
	var cache_time time.Time

	valkey := repostiories.Valkey
	cached_balance_as_bytes, err := valkey.Do(context.Background(), valkey.B().Get().Key(account.Id.String()).Build()).AsBytes()
	if err == nil {
		var cached_balance model.AccountBalance
		err = json.Unmarshal(cached_balance_as_bytes, &cached_balance)
		if err == nil {
			balance = cached_balance.Balance
			cache_time = cached_balance.Date
		}
	} else {
		balance_snapshot, err := repostiories.AccountRepository.GetBalanceSnapshot(account.Id.String())
		if err == nil && balance_snapshot != nil {
			balance = balance_snapshot.Balance
			cache_time = balance_snapshot.Date
		}
	}

	var transactions *[]model.Transaction
	if balance.GreaterThan(decimal.Zero) {
		transactions, err = repostiories.TransactionRepository.GetTransactionsByDate(account.Id.String(), cache_time.UTC(), now)
	} else {
		transactions, err = repostiories.TransactionRepository.GetAllTransactionsByAccount(account.Id.String())
	}

	if err == nil && transactions != nil {
		for _, transaction := range *transactions {
			switch transaction.Kind {
			case "deposit":
				balance = balance.Add(transaction.Amount)
			case "withdrawal":
				balance = balance.Sub(transaction.Amount)
			case "transfer":
				if transaction.FromAccountId != nil && *transaction.FromAccountId == account.Id {
					balance = balance.Sub(transaction.Amount)
				}

				if transaction.ToAccountId != nil && *transaction.ToAccountId == account.Id {
					balance = balance.Add(transaction.Amount)
				}
			case "refund":
				if transaction.FromAccountId != nil && *transaction.FromAccountId == account.Id {
					balance = balance.Add(transaction.Amount)
				}

				if transaction.ToAccountId != nil && *transaction.ToAccountId == account.Id {
					balance = balance.Sub(transaction.Amount)
				}
			default:
				return nil, errors.New("unknown transaction kind")
			}
		}
	}

	account_balance := model.AccountBalance{
		AccountId: account.Id,
		Balance:   balance,
		Date:      now,
	}

	b, err := json.Marshal(account_balance)
	if err == nil {
		valkey.Do(ctx, valkey.B().Del().Key(account.Id.String()).Build())
		valkey.Do(ctx, valkey.B().Set().Key(account.Id.String()).Value(string(b)).Ex(24*time.Hour).Build())
	}

	return &account_balance, nil
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
