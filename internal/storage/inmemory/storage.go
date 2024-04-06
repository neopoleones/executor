package inmemory

import (
	"context"
	"executor/internal/models"
	"executor/internal/storage"
	"github.com/google/uuid"
	"log/slog"
)

type DummyStorage struct {
	data map[uuid.UUID]*models.Runnable
}

func (d DummyStorage) GetCommands(ctx context.Context) ([]*models.Runnable, error) {
	result := make([]*models.Runnable, 0, len(d.data))

	for _, v := range d.data {
		result = append(result, v)
	}

	return result, nil
}

func (d DummyStorage) GetCommandByID(_ context.Context, uuid uuid.UUID) (*models.Runnable, error) {
	r, found := d.data[uuid]
	if found {
		return r, nil
	}
	return nil, storage.ErrNotFound
}

func (d DummyStorage) AddCommand(_ context.Context, sources []string) (*models.Runnable, error) {
	nr := models.NewRunnable(sources)
	d.data[nr.Sid] = nr

	return nr, nil
}

func (d DummyStorage) UpdateCommandInfo(ctx context.Context, runnable *models.Runnable) error {
	r, found := d.data[runnable.Sid]
	if !found {
		return storage.ErrNotFound
	}

	r.Status = runnable.Status
	r.Info = runnable.Info

	return nil
}

func (d DummyStorage) AddCommandOutput(_ context.Context, uuid uuid.UUID, output []string) error {
	r, found := d.data[uuid]
	if !found {
		return storage.ErrNotFound
	}

	r.Info.Output = append(r.Info.Output, output...)
	return nil
}

func (d DummyStorage) Close(ctx context.Context) {
	// Dummy: inmemory storage can't be closed
	slog.Info("Closed DummyStorage")
}

func GetStorage() (*DummyStorage, error) {
	return &DummyStorage{
		make(map[uuid.UUID]*models.Runnable),
	}, nil
}
