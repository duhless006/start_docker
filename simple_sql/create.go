package simple_sql

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func CreateTable(ctx context.Context, conn *pgx.Conn) error {
	sqlQuery := `
	CREATE TABLE IF NOT EXISTS personal (
		id SERIAL PRIMARY KEY,
		full_name VARCHAR(255) NOT NULL,
		position VARCHAR(255) NOT NULL
	);
	`
	_, err := conn.Exec(ctx, sqlQuery)
	return err
}
