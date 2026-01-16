package database

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func ConnectDataBase(ctx context.Context) (*pgx.Conn, error) {
	val := os.Getenv("CONN_STRING")
	return pgx.Connect(ctx, val)
}
