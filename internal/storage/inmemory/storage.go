package inmemory

import (
	"executor/internal/models"
	"github.com/google/uuid"
)

type DummyStorage struct {
	data map[uuid.UUID]*models.Runnable
}

func (d DummyStorage) GetCommands() ([]*models.Runnable, error) {
	//TODO implement me
	panic("implement me")
}

func (d DummyStorage) GetCommandByID(uuid uuid.UUID) (*models.Runnable, error) {
	//TODO implement me
	panic("implement me")
}

func (d DummyStorage) AddCommand(sources string) (*models.Runnable, error) {
	//TODO implement me
	panic("implement me")
}

func (d DummyStorage) UpdateCommandInfo(runnable *models.Runnable) error {
	//TODO implement me
	panic("implement me")
}

func GetStorage() (*DummyStorage, error) {
	return &DummyStorage{
		make(map[uuid.UUID]*models.Runnable),
	}, nil
}
