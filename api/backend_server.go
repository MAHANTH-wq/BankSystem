package api

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/token"
	"github.com/mahanth/simplebank/util"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// Function to create a new server instance
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	fmt.Println(err)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker for server")
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setUpRouter()

	return server, nil
}

// Function to start the http server on port address to start serving requests
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

//Function Set Up Server

func (server *Server) setUpRouter() {
	router := gin.Default()
	// Routes defined for the server

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccounts)
	authRoutes.PUT("/addbalance", server.addAccountBalance)
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router

}

// gin.H is a shortcut for map[string]interface{}
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
