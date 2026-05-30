package main

import (
	"SimpleBank/api"
	database "SimpleBank/database/sqlc"
	"SimpleBank/util"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	store := database.NewStore(conn)
	redisClient := util.NewRedisClient(config.RedisAddress)
	defer redisClient.Close()
	server, err := api.NewServer(config, store, redisClient)
	if err != nil {
		log.Fatal("cannot create server:", err)
	}
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server", err)
	}
}
