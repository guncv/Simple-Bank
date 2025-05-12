package main

import (
	"context"
	"database/sql"
	"os"

	"net"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/guncv/Simple-Bank/api"
	db "github.com/guncv/Simple-Bank/db/sqlc"
	_ "github.com/guncv/Simple-Bank/docs/statik"
	"github.com/guncv/Simple-Bank/gapi"
	"github.com/guncv/Simple-Bank/mail"
	pb "github.com/guncv/Simple-Bank/pb"
	"github.com/guncv/Simple-Bank/util"
	"github.com/guncv/Simple-Bank/worker"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db")
	}

	// run db migration
	runDBMigration(config.MigrationsURL, config.DBSource)

	store := db.NewStore(conn)

	log.Info().Msgf("redis address: %s", config.RedisAddress)
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProcessor(redisOpt, store, config)
	go runGatewayServer(taskDistributor, config, store)
	runGrpcServer(taskDistributor, config, store)
	// runGinServer(config, store)

}

func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatal().Msg("cannot create new migration")
	}

	if err := migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Msg("failed to up migration")
	}

	log.Info().Msg("db migrated successfully")
}

func runGrpcServer(taskDistributor worker.TaskDistributor, config util.Config, store db.Store) {
	server, err := gapi.NewServer(taskDistributor, config, store)
	if err != nil {
		log.Fatal().Msg("cannot create new server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Info().Msgf("start GRPC server at %s", listener.Addr().String())

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal().Msg("cannot start grpc server")
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.Store, config util.Config) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mailer)

	log.Info().Msg("task processor started")

	if err := taskProcessor.Start(); err != nil {
		log.Fatal().Err(err).Msg("fail to start task processor")
	}
}

func runGatewayServer(taskDistributor worker.TaskDistributor, config util.Config, store db.Store) {
	server, err := gapi.NewServer(taskDistributor, config, store)
	if err != nil {
		log.Fatal().Msg("cannot create new server")
	}
	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	grpcMux := runtime.NewServeMux(jsonOption)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("cannot register handler server")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Msg("cannot create statik fs")
	}
	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))
	mux.Handle("/swagger/", swaggerHandler)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot create listener")
	}

	log.Info().Msgf("start HTTP gateway server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	err = http.Serve(listener, handler)
	if err != nil {
		log.Fatal().Msg("cannot start HTTP gateway server")
	}
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create new server")
	}

	if err = server.Start(config.HTTPServerAddress); err != nil {
		log.Fatal().Msg("cannot start server")
	}

	log.Info().Msgf("start http server at %s", config.HTTPServerAddress)
}
