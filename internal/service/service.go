package service

import (
	"context"
	"executor/internal/storage"
)

type ExecutorService interface {
	Run(context.Context) error
	Setup(storage.ExecutorStorage)
	Release(context.Context)
}
