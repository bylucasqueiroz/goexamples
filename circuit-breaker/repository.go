package main

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
)

type Repository struct {
	Conn *pgx.Conn
}

func NewRepository(conn *pgx.Conn) *Repository {
	return &Repository{
		Conn: conn,
	}
}

func (db *Repository) Save(ctx context.Context, order Order) error {
	_, err := db.Conn.Exec(
		ctx,
		`INSERT INTO orders (name, amount) 
		VALUES ($1, $2)`,
		order.Name,
		order.Amount,
	)
	if err != nil {
		return fmt.Errorf("Error inserting event into database: %v", err)
	}

	return nil
}
