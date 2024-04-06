package executor

import (
	"context"
	"errors"
	"executor/internal/models"
	"github.com/google/uuid"
)

var (
	ErrNotScheduled = errors.New("this runnable isn't in scheduled state")
)

type CommandExecutor interface {
	// Run gets scheduled command from storage and executes it
	Run(context.Context, uuid.UUID) (*models.Runnable, error)
	Release(context.Context)

	// Start initializes scheduler
	Start(context.Context)
}
