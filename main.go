package main

import (
	"SimpleBank/api"
	database "SimpleBank/database/sqlc"
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	var err error
	err = godotenv.Load()
	conn, err := sql.Open(os.Getenv("DB_DRIVER"), os.Getenv("DB_SOURCE"))
	if err != nil {
		log.Fatal("cannot connect to database:", err)
	}
	store := database.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(os.Getenv("SERVER_ADDRESS"))
	if err != nil {
		log.Fatal("cannot start server")
	}
}
