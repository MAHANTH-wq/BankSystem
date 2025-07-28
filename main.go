package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/hibiken/asynq"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/util"
	"github.com/mahanth/simplebank/worker"
	"github.com/rakyll/statik/fs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/mahanth/simplebank/api"
	_ "github.com/mahanth/simplebank/doc/statik"
	"github.com/mahanth/simplebank/gapi"
	"github.com/mahanth/simplebank/pb"

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

	runDBMigration(config.MigrationURL, config.DBSource)

	store := db.NewStore(dbConnPool)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store)
	go runGrpcGateway(config, store, taskDistributor)
	runGrpcServer(config, store, taskDistributor)

}

func runDBMigration(migrationURL string, dbSource string) {
	migrationObject, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal("cannot create new migrate instance:", err)
	}

	if err = migrationObject.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("failed to run migrate up: ", err)
	}
	log.Println("db migrated successfully")
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Println("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal("error starting the task processor server %s", err)
	}
}
func runGrpcServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {

	server, err := gapi.NewServer(config, store, taskDistributor)

	if err != nil {
		fmt.Println(err)
		log.Fatal("error creating new server instance")
		panic(err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterBankSystemServer(grpcServer, server)

	reflection.Register(grpcServer)
	listener, err := net.Listen("tcp", config.GrpcServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start gRPC server")
	}

}

func runGrpcGateway(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {

	server, err := gapi.NewServer(config, store, taskDistributor)

	if err != nil {
		fmt.Println(err)
		log.Fatal("error creating new server instance")
		panic(err)
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcGatewayMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	err = pb.RegisterBankSystemHandlerServer(ctx, grpcGatewayMux, server)
	if err != nil {
		log.Fatal("cannot register handler server")
	}

	mux := http.NewServeMux()

	mux.Handle("/", grpcGatewayMux)

	//Swagger UI

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal("cannot create statik fs: ", err)
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HttpServerAddress)
	if err != nil {
		log.Fatal("cannot create listener")
	}
	log.Printf("start Http gateway server at %s", listener.Addr().String())
	err = http.Serve(listener, mux)

	if err != nil {
		log.Fatal("cannot start gRPC server")
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)

	if err != nil {
		fmt.Println(err)
		log.Fatal("error creating new server instance")
		panic(err)
	}

	err = server.Start(config.HttpServerAddress)
	if err != nil {
		log.Fatalf("error starting the server at port %s", config.HttpServerAddress)
		panic(err)
	}

}
