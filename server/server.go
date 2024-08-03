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
	router := gin.Default()

	server := Server{
		Repositories: repositories,
		Router:       router,
	}

	return &server
}

func (s *Server) Start(address string) {
	s.Router.Run(address)
}
