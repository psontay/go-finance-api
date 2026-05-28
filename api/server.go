package api

import (
	database "SimpleBank/database/sqlc"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)
import "github.com/gin-gonic/gin"

// server serves http request
type Server struct {
	store  database.Store
	router *gin.Engine
}

func NewServer(store database.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)

		if err != nil {
			return server
		}
	}

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)
	router.GET("/transfers/:id", server.getTransfer)
	router.GET("/transfers", server.listTransfers)

	server.router = router
	return server
}

// run http server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
