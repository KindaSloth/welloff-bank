package main

import (
	"bytes"
	"context"
	"net/http"
	"sync"
	"testing"
	"welloff-bank/server"
	"welloff-bank/utils"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

// TODO: create a test postgres and valkey instance

func TransferTransactionRequest(from_account_id string, to_account_id string, amount string) error {
	url := "http://localhost:5001/transaction/transfer"
	payload := []byte(`{
		"amount": "` + amount + `",
		"from_account_id": "` + from_account_id + `",
		"to_account_id": "` + to_account_id + `"
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	// TODO: remove this mock session id
	req.Header.Set("Cookie", "sessionId=01912a50-18ac-7859-944f-cb4745d34c7e; Max-Age=86400; Domain=localhost; Path=/; Secure; HttpOnly")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func DepositTransactionRequest(to_account_id string, amount string) error {
	url := "http://localhost:5001/transaction/deposit"
	payload := []byte(`{
		"amount": "` + amount + `",
		"to_account_id": "` + to_account_id + `"
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	// TODO: remove this mock session id
	req.Header.Set("Cookie", "sessionId=01912a50-18ac-7859-944f-cb4745d34c7e; Max-Age=86400; Domain=localhost; Path=/; Secure; HttpOnly")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func WithdrawalTransactionRequest(from_account_id string, amount string) error {
	url := "http://localhost:5001/transaction/withdrawal"
	payload := []byte(`{
		"amount": "` + amount + `",
		"from_account_id": "` + from_account_id + `"
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	// TODO: remove this mock session id
	req.Header.Set("Cookie", "sessionId=01912a50-18ac-7859-944f-cb4745d34c7e; Max-Age=86400; Domain=localhost; Path=/; Secure; HttpOnly")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func TestBalanceCorrectnessAfterMultipleTransactions(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Error(err)
	}
	server := server.New()

	go func() { server.Start("localhost:5001") }()

	DepositTransactionRequest("01911f41-f631-734b-a631-46aca9614536", "1000.00")

	account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse("01911f41-f631-734b-a631-46aca9614536"), server.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	numRequests := 1000
	parallelism := 100

	var wg sync.WaitGroup
	sem := make(chan struct{}, parallelism)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer func() { <-sem }()
			defer wg.Done()
			err := TransferTransactionRequest("01911f41-f631-734b-a631-46aca9614536", "01911fb7-727e-74e2-97d3-639eab3966ef", "1.00")
			if err != nil {
				t.Error(err)
			}
		}()
	}

	wg.Wait()

	new_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse("01911f41-f631-734b-a631-46aca9614536"), server.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	if account_balance.Balance.Sub(decimal.NewFromInt(1000)).String() != new_account_balance.Balance.String() {
		t.Errorf("Balance is not correct. Expected: %s, Actual: %s", account_balance.Balance.Sub(decimal.NewFromInt(1000)).String(), new_account_balance.Balance.String())
	}
}
