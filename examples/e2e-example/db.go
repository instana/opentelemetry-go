package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/XSAM/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.41.0"

	_ "modernc.org/sqlite" // pure-Go SQLite driver, no cgo required
)

// sqliteDriverName is the driver name registered by modernc.org/sqlite.
const sqliteDriverName = "sqlite"

// Order is the domain object returned by the order service.
type Order struct {
	ID       int64
	Customer string
	Item     string
	Amount   float64
}

// initDB opens an in-memory SQLite database instrumented with otelsql, creates
// the schema, and seeds a few rows.
//
// otelsql wraps the standard database/sql driver so that every Query/Exec
// automatically produces a child span under whatever context.Context is passed
// to the sql.DB methods. The span carries db.system, db.statement, and other
// semantic-convention attributes populated by the XSAM/otelsql library.
func initDB() (*sql.DB, error) {
	// otelsql.Open is a drop-in replacement for sql.Open. It accepts the same
	// functional options as sql.Open plus OTel-specific ones. WithAttributes
	// stamps every span with the db.system.name semantic convention attribute.
	// WithSpanNameFormatter controls the span name so it reads "db: <SQL>".
	db, err := otelsql.Open(
		sqliteDriverName,
		"file::memory:?cache=shared",
		otelsql.WithAttributes(semconv.DBSystemNameSQLite),
		otelsql.WithSpanNameFormatter(func(_ context.Context, _ otelsql.Method, query string) string {
			return "db: " + query
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("otelsql.Open: %w", err)
	}

	if err := createSchema(db); err != nil {
		db.Close()
		return nil, err
	}
	if err := seedData(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createSchema(db *sql.DB) error {
	_, err := db.ExecContext(context.Background(), `
		CREATE TABLE IF NOT EXISTS orders (
			id       INTEGER PRIMARY KEY,
			customer TEXT    NOT NULL,
			item     TEXT    NOT NULL,
			amount   REAL    NOT NULL
		)
	`)
	return err
}

func seedData(db *sql.DB) error {
	rows := []Order{
		{ID: 1, Customer: "Alice", Item: "Widget A", Amount: 29.99},
		{ID: 42, Customer: "Bob", Item: "Gadget B", Amount: 149.00},
		{ID: 100, Customer: "Carol", Item: "Doohickey C", Amount: 9.99},
	}
	for _, o := range rows {
		_, err := db.ExecContext(context.Background(),
			`INSERT INTO orders (id, customer, item, amount) VALUES (?, ?, ?, ?)`,
			o.ID, o.Customer, o.Item, o.Amount,
		)
		if err != nil {
			return fmt.Errorf("seed order %d: %w", o.ID, err)
		}
	}
	return nil
}

// queryOrder fetches a single order by ID.
//
// The ctx here carries the active server span created by otelhttp, so the
// child span that otelsql creates will be nested under it in the trace.
func queryOrder(ctx context.Context, db *sql.DB, id int64) (*Order, error) {
	row := db.QueryRowContext(ctx,
		`SELECT id, customer, item, amount FROM orders WHERE id = ?`, id,
	)

	var o Order
	if err := row.Scan(&o.ID, &o.Customer, &o.Item, &o.Amount); err != nil {
		return nil, err
	}
	return &o, nil
}
