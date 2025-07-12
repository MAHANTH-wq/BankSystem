package api

import (
	"github.com/go-playground/validator/v10"
	db "github.com/mahanth/simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// Function to create a new server instance
func NewServer(store db.Store) *Server {
	server := &Server{
		store: store,
	}
	router := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	// Routes defined for the server
	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)
	router.PUT("/addbalance", server.addAccountBalance)
	router.POST("/transfers", server.createTransfer)

	server.router = router
	return server
}

// Function to start the http server on port address to start serving requests
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// gin.H is a shortcut for map[string]interface{}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
