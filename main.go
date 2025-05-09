package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/guncv/Simple-Bank/api"
	db "github.com/guncv/Simple-Bank/db/sqlc"
	"github.com/guncv/Simple-Bank/gapi"
	pb "github.com/guncv/Simple-Bank/pb"
	"github.com/guncv/Simple-Bank/util"
	_ "github.com/lib/pq"
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
		log.Fatal("cannot connect to dv", err)
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
	runGinServer(config, store)

}

func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create new server:", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("cannot create listener:", err)
	}

	log.Printf("start grpc server at %s", listener.Addr().String())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("cannot start grpc server:", err)
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create new server:", err)
	}

	if err = server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}

	log.Printf("start http server at %s", config.HTTPServerAddress)
}
