package server

import (
	"welloff-bank/repository"

	"github.com/gin-gonic/gin"
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

	router.Use(s.AuthMiddleware())

	return router
}

func (s *Server) Start(addr string) {
	router := s.SetupRouter(addr)

	router.Run(addr)
}
