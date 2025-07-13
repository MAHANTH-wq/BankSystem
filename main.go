package main

import (
	"context"
	"fmt"
	"log"

	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/util"

	"github.com/mahanth/simplebank/api"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	var err error
	config, err := util.LoadConfig(".")

	if err != nil {
		log.Fatal("error loading configuration ")
		panic(err)
	}

	ctx := context.Background()
	dbConnPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatalf("error getting new db connection pool")
		panic(err)
	}

	store := db.NewStore(dbConnPool)

	server, err := api.NewServer(config, store)

	if err != nil {
		fmt.Println(err)
		log.Fatal("error creating new server instance")
		panic(err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatalf("error starting the server at port %s", config.ServerAddress)
		panic(err)
	}

}
