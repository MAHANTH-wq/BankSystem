package main

import (
	"context"
	"simplebank/api"
	db "simplebank/db/sqlc"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	httpServerAddress = "0.0.0.0:8080"
	dbSource          = "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	var err error
	ctx := context.Background()
	dbConnPool, err := pgxpool.New(ctx, dbSource)
	if err != nil {
		panic(err)
	}

	store := db.NewStore(dbConnPool)

	server := api.NewServer(store)

	err = server.Start(httpServerAddress)
	if err != nil {
		panic(err)
	}

}
