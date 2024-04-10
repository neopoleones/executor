package rest

import (
	"encoding/json"
	"errors"
	"executor/internal/models"
	"executor/internal/storage"
	"github.com/google/uuid"
	"net/http"
)

type KillRequest struct {
	RawSid string `json:"sid"`
}

type KillResponse struct {
	Status string `json:"status"`
}

func (s *Service) killHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var kr KillRequest

		// Parse request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&kr); len(kr.RawSid) == 0 {
			if err == nil {
				err = errors.New("no sid provided")
			} else {
				err = errors.New("bad request")
			}

			response(NewErrorResponse(err), http.StatusBadRequest, w)
			return
		}

		// Check if sid is correct
		sid, err := uuid.Parse(kr.RawSid)
		if err != nil {
			response(NewErrorResponse(errors.New("bad sid")), http.StatusBadRequest, w)
			return
		}

		// Get runnable
		runnable, err := s.storage.GetCommandByID(r.Context(), sid)
		if err != nil {
			var code = http.StatusInternalServerError
			if errors.Is(err, storage.ErrNotFound) {
				code = http.StatusNotFound
			}

			response(NewErrorResponse(err), code, w)
			return
		}

		// Is execution in progress?
		if runnable.Status != models.StatusInProgress {
			response(NewErrorResponse(errors.New("command not in running state")), http.StatusUnprocessableEntity, w)
			return
		}

		// Update status: scheduler will kill the runnable by itself
		runnable.UpdateStatus(models.StatusRejected)
		if err := s.storage.UpdateCommandInfo(r.Context(), runnable); err != nil {
			response(NewErrorResponse(err), http.StatusInternalServerError, w)
		} else {
			response(KillResponse{
				Status: "ok",
			}, http.StatusOK, w)
		}
	}
}
