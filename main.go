package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	db "github.com/mahanth/simplebank/db/sqlc"
	"github.com/mahanth/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/mahanth/simplebank/api"
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

	store := db.NewStore(dbConnPool)

	go runGrpcGateway(config, store)
	runGrpcServer(config, store)

}

func runGrpcServer(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)

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

func runGrpcGateway(config util.Config, store db.Store) {

	server, err := gapi.NewServer(config, store)

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
