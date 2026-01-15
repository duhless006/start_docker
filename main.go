package main

import (
	"context"
	database "copr/Database"
	httppack "copr/httpPack"
	"copr/simple_sql"
	"copr/worker"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

func main() {

	ctx := context.Background()

	log.Println("Starting application...")
	log.Printf("Database URL: %s", os.Getenv("CONN_STRING"))

	var conn *pgx.Conn
	var err error

	for i := 0; i < 5; i++ {
		conn, err = database.ConnectDataBase(ctx)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err := conn.Ping(ctx); err != nil {
		log.Fatal("Database ping failed:", err)
	}
	defer func() {
		conn.Close(ctx)
		log.Println("Database connection closed")
	}()

	log.Println("Database connection established")

	if err := simple_sql.CreateTable(ctx, conn); err != nil {
		log.Fatal("Failed to create table:", err)
	}
	log.Println("Table 'personal' created/verified")

	log.Println("Server starting...")

	workers := worker.NewWorkers()
	httpHandlers := httppack.NewHTTPHandlers(workers, conn)
	httpServer := httppack.NewHTTPServer(*httpHandlers)

	if err := httpServer.ConnectServer(); err != nil {
		log.Fatal("Failed to start HTTP server:", err)
	}
}
