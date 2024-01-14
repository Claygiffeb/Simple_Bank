package main

import (
	"database/sql"
	"log"

	api "github.com/Clayagiffeb/Simple_Bank/API"
	db "github.com/Clayagiffeb/Simple_Bank/db/sqlc"
	"github.com/Clayagiffeb/Simple_Bank/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Failed to load config")
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}
	// create database
	store := db.NewStore(conn)
	// create server
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Failed to start server: ", err)
	}

}
