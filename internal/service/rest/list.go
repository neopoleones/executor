package rest

import (
	"executor/internal/models"
	"github.com/google/uuid"
	"net/http"
)

func (s *Service) listHandler() http.HandlerFunc {
	type RunnableSimplified struct {
		Sid    uuid.UUID             `json:"sid"`
		Status models.RunnableStatus `json:"status"`
	}

	type ListResponse struct {
		Status   string               `json:"status"`
		Commands []RunnableSimplified `json:"commands"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		rl, err := s.storage.GetCommands(r.Context())
		if err != nil {
			response(NewErrorResponse(err), http.StatusInternalServerError, w)
			return
		}

		// Omit details
		rs := make([]RunnableSimplified, len(rl))
		for i, v := range rl {
			rs[i] = RunnableSimplified{
				Sid:    v.Sid,
				Status: v.Status,
			}
		}

		response(ListResponse{
			Status:   "ok",
			Commands: rs,
		}, http.StatusOK, w)
	}
}
