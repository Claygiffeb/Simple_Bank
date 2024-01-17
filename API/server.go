package api

import (
	"fmt"

	db "github.com/Clayagiffeb/Simple_Bank/db/sqlc"
	"github.com/Clayagiffeb/Simple_Bank/token"
	"github.com/Clayagiffeb/Simple_Bank/util"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server which serves the HTTP requests
type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker *token.JWTMaker
}

// NewServer creates a new Server and set up routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(util.RandomString(32))
	if err != nil {
		return nil, fmt.Errorf("Faild to create token: %v", err)
	}

	server := &Server{
		store:      store,
		tokenMaker: tokenMaker,
		config:     config,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users", server.CreateUser)      // API for creating users
	router.POST("/users/login", server.loginUser) // API for login user
	router.POST("tokens/renew_access", server.renewAccessToken)
	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.CreateAccount)   // API for creating accounts
	authRoutes.GET("/accounts/:id", server.GetAccount)   // API for getting accounts
	authRoutes.GET("/accounts", server.ListAccounts)     // API for listing accounts
	authRoutes.POST("/transfers", server.CreateTransfer) // API for creating transfer
	server.router = router
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
