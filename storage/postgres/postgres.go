package postgres

import (
	"context"
	"fmt"
	"os"
	"time"

	"messaggio/storage/sqlc/database"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var log = logrus.New()

// PSQLClient represents the PostgreSQL client with a connection pool and queries.
type PSQLClient struct {
	DB      *pgxpool.Pool
	Queries *database.Queries
}

// NewPSQLClient creates a new PSQLClient instance with the provided configuration.
func NewPSQLClient() *PSQLClient {
	return &PSQLClient{}
}

// Connect establishes a connection to the PostgreSQL database using the provided credentials.
// It sets up connection pooling with maximum open and idle connections,
// and sets the maximum lifetime of connections.
func (s *PSQLClient) Connect(user, password, host, port, dbname string) error {
	const op = "storage.postgres.Connect()"

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return fmt.Errorf("%s: failed to parse DSN: %w", op, err)
	}

	config.MaxConns = 750
	config.MinConns = 10

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return fmt.Errorf("%s: failed to create connection pool: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		db.Close()
		return fmt.Errorf("%s: failed to ping database: %w", op, err)
	}

	s.DB = db
	s.Queries = database.New(db)
	if s.Queries == nil {
		db.Close()
		return fmt.Errorf("%s: failed to initialize Queries", op)
	}

	return nil
}

// Close closes the connection to the PostgreSQL database.
// It checks if there is an active database connection (s.DB) and attempts to close it.
// Logs an error message if there was an issue closing the connection.
func (s *PSQLClient) Close() {
	if s.DB != nil {
		s.DB.Close()
		if os.Getenv("FIBER_PREFORK_CHILD") == "" {
			log.Info("Connection to PostgreSQL closed")
		}
	}
}
