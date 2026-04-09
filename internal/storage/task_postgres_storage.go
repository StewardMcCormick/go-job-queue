package storage

import (
	"context"

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
	// query := `SELECT * FROM`
	//
	return nil, pgx.ErrNoRows
}
