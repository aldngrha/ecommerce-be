package database

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
)

func ConnectionDB(ctx context.Context, connectionString string) *sql.DB {
	db, err := sql.Open("postgres", connectionString)

	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}

	err = db.PingContext(ctx)
	if err != nil {
		panic("Failed to connect to the database: " + err.Error())
	}
	return db
}
