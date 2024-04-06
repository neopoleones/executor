package rest

import (
	"errors"
	"executor/internal/models"
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
			response(NewErrorResponse(err), http.StatusBadRequest, w)
			return
		}

		response(GetResponse{
			Status:   "ok",
			Runnable: runnable,
		}, http.StatusOK, w)
	}
}
