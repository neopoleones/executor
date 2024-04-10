package rest

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type ScheduleRequest struct {
	Script []string `json:"script"`
}

type ScheduleResponse struct {
	Status    string    `json:"status"`
	Sid       string    `json:"sid"`
	SchedTime time.Time `json:"sched_time"`
}

func (s *Service) scheduleHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		var sr ScheduleRequest

		// Parse request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&sr); len(sr.Script) == 0 {
			if err == nil {
				err = errors.New("no script provided")
			} else {
				// shadow the error: small hardening
				err = errors.New("bad request")
			}

			response(NewErrorResponse(err), http.StatusBadRequest, w)
			return
		}

		// Add command to storage as scheduled
		// Executor is checking for the new commands on background and running all of them

		if cmd, err := s.storage.AddCommand(r.Context(), sr.Script); err != nil {
			response(NewErrorResponse(err), http.StatusInternalServerError, w)
		} else {
			resp := ScheduleResponse{
				Status:    "ok",
				Sid:       cmd.Sid.String(),
				SchedTime: cmd.Info.ScheduledTime,
			}

			response(resp, http.StatusOK, w)
		}
	}
}
