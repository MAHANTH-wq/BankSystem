package main

import (
	"context"

	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/util"

	"github.com/mahanth/simplebank/api"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var err error
	config, err := util.LoadConfig(".")

	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	dbConnPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		panic(err)
	}

	store := db.NewStore(dbConnPool)

	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		panic(err)
	}

}
