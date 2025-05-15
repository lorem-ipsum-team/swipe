package postgres

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
	postgres_mig "github.com/lorem-ipsum-team/swipe/db/postgres"
	"github.com/lorem-ipsum-team/swipe/pkg/logger"
)

type Repo struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func NewRepo(ctx context.Context, log *slog.Logger, connString string) (*Repo, error) {
	log = log.WithGroup("postgres_repo")
	log.Info("connect to db", slog.Any("connection string", logger.Secret(connString)))
	log.InfoContext(ctx, "running migrations")

	err := postgres_mig.Up(ctx, connString)
	if err != nil {
		return nil, fmt.Errorf("failed to run migration: %w", err)
	}

	log.InfoContext(ctx, "creating pgx connection pool")

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	return &Repo{
		pool: pool,
		log:  log,
	}, nil
}
