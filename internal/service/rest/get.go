package rest

import (
	"errors"
	"executor/internal/models"
	"executor/internal/storage"
	"github.com/google/uuid"
	"net/http"
)

func (s *Service) getHandler() http.HandlerFunc {
	const (
		paramSid = "sid"
	)

	type GetResponse struct {
		Status   string           `json:"status"`
		Runnable *models.Runnable `json:"runnable"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rawSid := r.URL.Query().Get(paramSid)
		sid, err := uuid.Parse(rawSid)

		if len(rawSid) == 0 || err != nil {
			response(NewErrorResponse(errors.New("sid not found or incorrect")), http.StatusBadRequest, w)
			return
		}

		runnable, err := s.storage.GetCommandByID(r.Context(), sid)
		if err != nil {
			var code = http.StatusInternalServerError
			if errors.Is(err, storage.ErrNotFound) {
				code = http.StatusNotFound
			}
			
			response(NewErrorResponse(err), code, w)
			return
		}

		response(GetResponse{
			Status:   "ok",
			Runnable: runnable,
		}, http.StatusOK, w)
	}
}
