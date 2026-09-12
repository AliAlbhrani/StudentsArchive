package engine

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AliAlbhrani/StudentsArchive/env"
	"github.com/AliAlbhrani/StudentsArchive/sqlc"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	PgxDB   *pgxpool.Pool
	Queries *sqlc.Queries
)

var _ = func() bool {
	ctx := context.TODO()

	config, err := pgxpool.ParseConfig(fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", env.DB_USER, env.DB_PASSWORD, env.DB_HOST, env.DB_PORT, env.DB_NAME))
	if err != nil {
		slog.Error("failed to parse database config", "error", err)
		return false
	}

	conn, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		return false
	}

	slog.Info("connected to database")

	PgxDB = conn
	Queries = sqlc.New(PgxDB)

	return true
}()

func IsUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code == "23505"
	}
	return false
}
