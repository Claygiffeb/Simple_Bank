package api

import (
	db "github.com/Clayagiffeb/Simple_Bank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Server which serves the HTTP requests
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new Server and set up routing
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.CreateUser)         // API for creating users
	router.POST("/accounts", server.CreateAccount)   // API for creating accounts
	router.GET("/accounts/:id", server.GetAccount)   // API for getting accounts
	router.GET("/accounts", server.ListAccounts)     // API for listing accounts
	router.POST("/transfers", server.CreateTransfer) // API for creating transfer

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
