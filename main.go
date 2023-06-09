package main

import (
	"database/sql"
	"log"
	"net"

	_ "github.com/lib/pq"
	"github.com/chau-doan/simplebank/api"
	db "github.com/chau-doan/simplebank/db/sqlc"
	"github.com/chau-doan/simplebank/gapi"
	"github.com/chau-doan/simplebank/pb"
	"github.com/chau-doan/simplebank/util"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	store := db.NewStore(conn)
	runGinServer(config, store)
}

func runGrpcServer(config util.Config, store db.Store){
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot connect to server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start gRPC server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil{
		log.Fatal("cannot create gRPC server")
	}
}

func runGinServer(config util.Config, store db.Store){
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot connect to server:", err)
	}

	err = server.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}

}