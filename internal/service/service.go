package service

import (
	"context"
	"executor/internal/executor"
	"executor/internal/storage"
)

type ExecutorService interface {
	Run(context.Context) error
	Setup(storage.ExecutorStorage, executor.CommandExecutor)
	Release(context.Context)
}
