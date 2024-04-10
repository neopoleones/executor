package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"executor/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

const scheduleEndpoint = "/cmd/schedule"

func scheduleRequest(h http.Handler, rd ScheduleRequest) (*httptest.ResponseRecorder, error) {
	reqData, err := json.Marshal(rd)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, scheduleEndpoint, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	return w, nil
}

func TestScheduleHandlerOK(t *testing.T) {
	// Test request data are correct
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		testCtx := context.Background()

		w := testPerfromRequestWithData[ScheduleRequest](t, h, scheduleRequest, ScheduleRequest{[]string{
			"uname -a",
			"id",
		}})
		assertCode(t, w, http.StatusOK)

		// Check output format
		resp := testUnmarshal[ScheduleResponse](t, w.Body)
		assertStatus(t, resp.Status, "ok")

		// Check that uuid was returned
		sid := testParseUUID(t, resp.Sid)

		// uuid should be correct
		if _, err := es.GetCommandByID(testCtx, sid); err != nil {
			t.Fatalf("failed to get runnable by sid: %v", err)
		}
	})
}

func TestScheduleHandlerBadRequest(t *testing.T) {
	// Empty script not allowed
	performTest(func(es storage.ExecutorStorage, h http.Handler) {

		w := testPerfromRequestWithData[ScheduleRequest](t, h, scheduleRequest, ScheduleRequest{[]string{}})
		assertCode(t, w, http.StatusBadRequest)

		resp := testUnmarshal[ScheduleResponse](t, w.Body)
		assertStatus(t, resp.Status, "err")

		// Check that commands not added
		cmds, _ := es.GetCommands(context.Background())
		if len(cmds) != 0 {
			t.Fatalf("command was created but it should be abandoned: %v", len(cmds))
		}
	})
}
