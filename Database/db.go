package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func ConnectDataBase(ctx context.Context) (*pgx.Conn, error) {
	err := godotenv.Load("Error loading .env file")
	if err != nil {
		log.Fatal()
	}
	val := os.Getenv("CONN_STRING")
	return pgx.Connect(ctx, val)
}
