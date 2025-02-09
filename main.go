package main

import (
	"database/sql"
	"log"

	"github.com/guncv/Simple-Bank/api"
	db "github.com/guncv/Simple-Bank/db/sqlc"
	"github.com/guncv/Simple-Bank/util"
	_ "github.com/lib/pq"
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
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create new server:", err)
	}

	if err = server.Start(config.ServerAddress); err != nil {
		log.Fatal("cannot start server:", err)
	}
}
