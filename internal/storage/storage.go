package storage

import (
	"executor/internal/models"
	"github.com/google/uuid"
)

type ExecutorStorage interface {
	GetCommands() ([]*models.Runnable, error)
	GetCommandByID(uuid.UUID) (*models.Runnable, error)

	AddCommand(sources string) (*models.Runnable, error)
	UpdateCommandInfo(*models.Runnable) error
}
