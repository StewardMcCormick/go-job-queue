package storage

import (
	"context"
	"fmt"

	pb "github.com/StewardMcCormick/go-job-queue/gen/go/api/v1"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskPostgresStorage struct {
	pool *pgxpool.Pool
}

func NewTaskPostgresStorage(pool *pgxpool.Pool) *taskPostgresStorage {
	return &taskPostgresStorage{
		pool: pool,
	}
}

func (s *taskPostgresStorage) GetById(ctx context.Context, id string) ([]*pb.Task, error) {
	query := `SELECT * FROM tasks WHERE id=$1`

	_, err := s.pool.Exec(ctx, query, id)
	if err != nil {
		return nil, fmt.Errorf("get task by id from postgres error: %w", err)
	}

	return nil, pgx.ErrNoRows
}
