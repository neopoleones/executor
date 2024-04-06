package storage

import "executor/internal/models"

func FilterRunnablesByStatus(commands []*models.Runnable, status models.RunnableStatus) []*models.Runnable {
	rl := make([]*models.Runnable, 0)

	for _, cmd := range commands {
		if cmd.Status == status {
			rl = append(rl, cmd)
		}
	}

	return rl
}
