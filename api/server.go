package api

import (
	database "SimpleBank/database/sqlc"
	"SimpleBank/token"
	"SimpleBank/util"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/redis/go-redis/v9"
)
import "github.com/gin-gonic/gin"

// server serves http request
type Server struct {
	store      database.Store
	router     *gin.Engine
	config     util.Config
	tokenMaker token.Maker
	redis      *redis.Client
}

func NewServer(config util.Config, store database.Store, redisClient *redis.Client) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
		redis:      redisClient,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("currency", validCurrency)

		if err != nil {
			return server, fmt.Errorf("cannot register currency validator: %w", err)
		}
		err = v.RegisterValidation("role", validRole)

		if err != nil {
			return server, fmt.Errorf("cannot register role validator: %w", err)
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

	// ROUTES
	authRoutes := router.Group("/")
	authRoutes.Use(authMiddleware(server.tokenMaker))
	clientRoutes := authRoutes.Group("/")
	clientRoutes.Use(roleMiddleware("depositor", "admin"))
	adminRoutes := authRoutes.Group("/")
	adminRoutes.Use(roleMiddleware("admin"))

	// ADMIN ROUTES
	adminRoutes.POST("/users", server.createUser)
	adminRoutes.GET("/transfers", server.listTransfers)
	adminRoutes.GET("/users/:username", server.getUser)
	adminRoutes.GET("/users", server.listUsers)
	adminRoutes.PUT("/users", server.updateUser)
	adminRoutes.DELETE("/users/:username", server.deleteUser)
	adminRoutes.GET("/accounts/admin/list", server.listAccounts)
	// auth
	router.POST("/users/login", server.loginUser)
	router.POST("/tokens/renew_access", server.renewAccessToken)

	// CLIENT ROUTES
	clientRoutes.POST("/accounts", server.createAccount)
	clientRoutes.GET("/accounts/:id", server.getAccount)
	clientRoutes.GET("/accounts", server.listAccountsOwner)
	clientRoutes.POST("/transfers", server.createTransfer)
	clientRoutes.GET("/transfers/:id", server.getTransfer)

	// PUBLIC ROUTES
	router.POST("/users/register", server.registerUser)

	server.router = router

}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
