package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"testing"
	"time"
	"welloff-bank/server"
	"welloff-bank/utils"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/shopspring/decimal"
)

// TODO: create a test postgres and valkey instance
var s *server.Server
var sessionId string

// TODO: remove these mock accounts
var user_account = "01911f41-f631-734b-a631-46aca9614536"
var transferable_account = "01911fb7-727e-74e2-97d3-639eab3966ef"

func TestMain(m *testing.M) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	s = server.New()
	go s.Start("localhost:5001")

	time.Sleep(time.Second * 2)

	login_url := "http://localhost:5001/login"
	payload := []byte(`{
		"email": "test@email.com",
		"password": "test123"
	}`)

	req, err := http.NewRequest("POST", login_url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	cookies := resp.Header["Set-Cookie"]
	if len(cookies) == 0 {
		log.Fatal("Missing cookie")
	}

	sessionId = strings.Split(strings.Split(cookies[0], ";")[0], "=")[1]

	os.Exit(m.Run())
}

func TransferTransactionRequest(from_account_id string, to_account_id string, amount string) (string, error) {
	url := "http://localhost:5001/transaction/transfer"
	payload := []byte(`{
		"amount": "` + amount + `",
		"from_account_id": "` + from_account_id + `",
		"to_account_id": "` + to_account_id + `"
	}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "sessionId="+sessionId+"; Max-Age=86400; Domain=localhost; Path=/; Secure; HttpOnly")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	var payload_resp interface{}
	err = json.NewDecoder(resp.Body).Decode(&payload_resp)
	if err != nil {
		return "", err
	}

	transaction_id, ok := payload_resp.(map[string]interface{})["payload"].(map[string]interface{})["transaction_id"].(string)
	if !ok {
		return "", fmt.Errorf("failed to parse transaction id")
	}

	return transaction_id, nil
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
	req.Header.Set("Cookie", "sessionId="+sessionId+"; Max-Age=86400; Domain=localhost; Path=/; Secure; HttpOnly")

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
	req.Header.Set("Cookie", "sessionId="+sessionId+"; Max-Age=86400; Domain=localhost; Path=/; Secure; HttpOnly")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func RefundTransactionRequest(transaction_id string) error {
	url := "http://localhost:5001/transaction/refund/" + transaction_id
	payload := []byte(`{}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Cookie", "sessionId="+sessionId+"; Max-Age=86400; Domain=localhost; Path=/; Secure; HttpOnly")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	return nil
}

func TestBalanceCorrectnessAfterMultipleDepositAndWithdrawalTransactions(t *testing.T) {
	num_deposits := 50
	num_withdrawals := 25
	parallelism := 10

	user_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse(user_account), s.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, parallelism)

	for i := 0; i < num_deposits; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer func() { <-sem }()
			defer wg.Done()
			err := DepositTransactionRequest(user_account, "1.00")
			if err != nil {
				t.Error(err)
			}
		}()
	}

	for i := 0; i < num_withdrawals; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer func() { <-sem }()
			defer wg.Done()
			err := WithdrawalTransactionRequest(user_account, "1.00")
			if err != nil {
				t.Error(err)
			}
		}()
	}

	wg.Wait()

	new_user_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse(user_account), s.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	expected_user_account_balance := user_account_balance.Balance.Add(decimal.NewFromInt(int64(num_deposits - num_withdrawals))).String()
	if expected_user_account_balance != new_user_account_balance.Balance.String() {
		t.Errorf("User Account Balance is not correct. Expected: %s, Actual: %s", expected_user_account_balance, new_user_account_balance.Balance.String())
	}
}

func TestBalanceCorrectnessAfterMultipleTransactions(t *testing.T) {
	DepositTransactionRequest(user_account, "100.00")

	user_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse(user_account), s.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	transferable_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse(transferable_account), s.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	requests_num := 100
	parallelism := 10

	var wg sync.WaitGroup
	sem := make(chan struct{}, parallelism)

	for i := 0; i < requests_num; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func() {
			defer func() { <-sem }()
			defer wg.Done()
			_, err := TransferTransactionRequest(user_account, transferable_account, "1.00")
			if err != nil {
				t.Error(err)
			}
		}()
	}

	wg.Wait()

	new_user_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse(user_account), s.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	new_transferable_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse(transferable_account), s.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	expected_user_account_balance := user_account_balance.Balance.Sub(decimal.NewFromInt(int64(requests_num))).String()
	if expected_user_account_balance != new_user_account_balance.Balance.String() {
		t.Errorf("User Account Balance is not correct. Expected: %s, Actual: %s", expected_user_account_balance, new_user_account_balance.Balance.String())
	}

	expected_transferable_account_balance := transferable_account_balance.Balance.Add(decimal.NewFromInt(int64(requests_num))).String()
	if expected_transferable_account_balance != new_transferable_account_balance.Balance.String() {
		t.Errorf("Transferable Account Balance is not correct. Expected: %s, Actual: %s", expected_transferable_account_balance, new_transferable_account_balance.Balance.String())
	}
}

func TestBalanceCorrectnessAfterMultipleTransactionsWithRefund(t *testing.T) {
	DepositTransactionRequest(user_account, "100.00")

	user_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse(user_account), s.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	transferable_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse(transferable_account), s.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	requests_num := 100
	refund_num := 50
	parallelism := 10

	var wg sync.WaitGroup
	sem := make(chan struct{}, parallelism)

	for i := 0; i < requests_num; i++ {
		wg.Add(1)
		sem <- struct{}{}

		go func(i int) {
			defer func() { <-sem }()
			defer wg.Done()
			transaction_id, err := TransferTransactionRequest(user_account, transferable_account, "1.00")
			if err != nil {
				t.Error(err)
			}
			if i+1 <= refund_num {
				err = RefundTransactionRequest(transaction_id)
				if err != nil {
					t.Error(err)
				}
			}
		}(i)
	}

	wg.Wait()

	new_user_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse(user_account), s.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	new_transferable_account_balance, err := utils.GetAccountBalance(context.Background(), uuid.MustParse(transferable_account), s.Repositories, false)
	if err != nil {
		t.Error(err)
	}

	expected_user_account_balance := user_account_balance.Balance.Sub(decimal.NewFromInt(int64(requests_num - refund_num))).String()
	if expected_user_account_balance != new_user_account_balance.Balance.String() {
		t.Errorf("User Account Balance is not correct. Expected: %s, Actual: %s", expected_user_account_balance, new_user_account_balance.Balance.String())
	}

	expected_transferable_account_balance := transferable_account_balance.Balance.Add(decimal.NewFromInt(int64(requests_num - refund_num))).String()
	if expected_transferable_account_balance != new_transferable_account_balance.Balance.String() {
		t.Errorf("Transferable Account Balance is not correct. Expected: %s, Actual: %s", expected_transferable_account_balance, new_transferable_account_balance.Balance.String())
	}
}
