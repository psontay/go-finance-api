package api

import (
	database "SimpleBank/database/sqlc"
	"SimpleBank/token"
	"SimpleBank/util"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)
import "github.com/gin-gonic/gin"

// server serves http request
type Server struct {
	store      database.Store
	router     *gin.Engine
	config     util.Config
	tokenMaker token.Maker
}

func NewServer(config util.Config, store database.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)

		if err != nil {
			return server, fmt.Errorf("cannot binding validator: %w", err)
		}
	}

	server.setupRouter()

	return server, nil
}

// run http server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) setupRouter() {
	router := gin.Default()
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	// PRIVATE ROUTES
	// account
	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	// transfer
	authRoutes.POST("/transfers", server.createTransfer)
	authRoutes.GET("/transfers/:id", server.getTransfer)
	authRoutes.GET("/transfers", server.listTransfers)
	// user
	authRoutes.GET("/users/:username", server.getUser)
	authRoutes.GET("/users", server.listUsers)
	authRoutes.PUT("/users", server.updateUser)
	authRoutes.DELETE("/users/:username", server.deleteUser)

	// PUBLIC ROUTES
	// user
	router.POST("/users", server.createUser)
	// auth
	router.POST("/users/login", server.loginUser)
	server.router = router

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
