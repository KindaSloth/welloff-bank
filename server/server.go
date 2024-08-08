package server

import (
	"context"
	"log"
	"welloff-bank/model"
	"welloff-bank/repository"
	"welloff-bank/utils"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

type Server struct {
	Repositories repository.Repositories
	Router       *gin.Engine
}

func New() *Server {
	repositories := repository.New()

	server := Server{
		Repositories: repositories,
	}

	return &server
}

func (s *Server) SetupRouter(addr string) *gin.Engine {
	router := gin.Default()
	router.Use(CorsMiddleware())

	router.GET("/health-check", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"message": "OK"})
	})
	router.POST("/register", s.Register())
	router.POST("/login", s.Login())

	router.Use(s.AuthMiddleware())

	// User enpoints
	router.GET("/me", s.Me())

	// Account enpoints
	router.POST("/account", s.CreateAccount())
	router.GET("/account/:id", s.GetAccount())
	router.GET("/accounts", s.GetAccounts())
	router.DELETE("/account/:id", s.DisableAccount())

	// Transaction enpoints
	router.GET("/transaction/:id", s.GetTransaction())
	router.POST("/transaction/deposit", s.DepositTransaction())
	router.POST("/transaction/withdrawal", s.WithdrawalTransaction())
	router.POST("/transaction/transfer", s.TransferTransaction())
	router.POST("/transaction/refund/:id", s.RefundTransaction())

	return router
}

func (s *Server) StartCron() {
	c := cron.New()
	c.AddFunc("@every 10m", func() {
		log.Println("[INFO] [Balance Snapshot Updater] running...")

		limit := 100
		offset := 0
		up_to_date := false

		for !up_to_date {
			accounts, err := s.Repositories.AccountRepository.GetAccounts(limit, offset)
			if err != nil {
				log.Println("[ERROR] [Balance Snapshot Updater] failed to get accounts: ", err)
				return
			}

			balances := make([]model.AccountBalance, 0)
			for _, account := range *accounts {
				balance, err := utils.GetAccountBalance(context.Background(), account.Id, s.Repositories, false)
				if err != nil {
					log.Println("[ERROR] [Balance Snapshot Updater] failed to get account balance: ", err)
					return
				}

				balances = append(balances, *balance)
			}

			err = s.Repositories.AccountRepository.BulkUpsertBalanceSnapshots(&balances)
			if err != nil {
				log.Println("[ERROR] [Balance Snapshot Updater] failed to upsert balance snapshots: ", err)
				return
			}

			offset += limit
			if len(*accounts) < limit {
				up_to_date = true
			}
		}

		log.Println("[INFO] [Balance Snapshot Updater] completed")
	})
	c.Start()
}

func (s *Server) Start(addr string) {
	router := s.SetupRouter(addr)

	router.Run(addr)
}
