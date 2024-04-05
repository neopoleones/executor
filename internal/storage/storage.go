package storage

import (
	"context"
	"errors"
	"executor/internal/models"
	"github.com/google/uuid"
)

var (
	ErrNotFound = errors.New("runnable not found")
)

type ExecutorStorage interface {
	GetCommands(context.Context) ([]*models.Runnable, error)
	GetCommandByID(context.Context, uuid.UUID) (*models.Runnable, error)

	AddCommand(context.Context, []string) (*models.Runnable, error)

	UpdateCommandInfo(context.Context, *models.Runnable) error
	AddCommandOutput(context.Context, uuid.UUID, []string) error

	Close(context.Context)
}
