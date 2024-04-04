package service

import (
	"context"
	"executor/internal/storage"
)

type ExecutorService interface {
	Run(ctx context.Context) error
	Setup(storage.ExecutorStorage)
}
