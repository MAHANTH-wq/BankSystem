package gapi

import (
	"fmt"

	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/pb"
	"github.com/mahanth/simplebank/token"
	"github.com/mahanth/simplebank/util"
)

type Server struct {
	pb.UnimplementedBankSystemServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// Function to create a new gRPC server instance
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker for server")
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	return server, nil
}
