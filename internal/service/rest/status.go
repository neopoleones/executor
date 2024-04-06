package rest

import (
	"executor/internal/models"
	"executor/internal/storage"
	"net/http"
	"time"
)

func (s *Service) statusHandler() http.HandlerFunc {
	type RunnableInfo struct {
		Running int `json:"running"`
		Exited  int `json:"exited"`
	}

	type StatusResponse struct {
		Status  string        `json:"status"`
		Version string        `json:"version"`
		Uptime  time.Duration `json:"uptime"`

		Commands RunnableInfo `json:"commands"`
	}

	handleInitialized := time.Now()

	return func(w http.ResponseWriter, r *http.Request) {
		// Get from storage all Runnable entities
		entities, err := s.storage.GetCommands(r.Context())
		if err != nil {
			response(NewErrorResponse(err), http.StatusInternalServerError, w)
			return
		}

		resp := StatusResponse{
			Status:  "ok",
			Version: s.version,
			Uptime:  time.Now().Sub(handleInitialized),
			Commands: RunnableInfo{
				Running: len(storage.FilterRunnablesByStatus(entities, models.StatusInProgress)),
				Exited:  len(storage.FilterRunnablesByStatus(entities, models.StatusDone)),
			},
		}

		response(resp, http.StatusOK, w)
	}
}
