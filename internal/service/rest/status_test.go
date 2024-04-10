package rest

import (
	"executor/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const statusEndpoint = "/status"

func statusRequest(h http.Handler) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(http.MethodGet, statusEndpoint, nil)
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	return w, nil
}

func TestStatusHandlerNoCommands(t *testing.T) {
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		w := testPerformRequest(t, h, statusRequest)
		assertCode(t, w, http.StatusOK)

		resp := testUnmarshal[StatusResponse](t, w.Body)

		if resp.Status != "ok" || resp.Commands.Exited != 0 || resp.Commands.Running != 0 {
			t.Fatalf("incorrect data returned: %+v", resp)
		}
	})
}

func TestStatusHandlerHasRunnable(t *testing.T) {
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		t.Parallel()

		// Create at least one runnable
		// We should wait for a small delay. For scheduler to take our runnable
		_ = testPerfromRequestWithData[ScheduleRequest](t, h, scheduleRequest, ScheduleRequest{[]string{
			"id",
		}})
		time.Sleep(time.Millisecond * 300)

		w := testPerformRequest(t, h, statusRequest)
		assertCode(t, w, http.StatusOK)

		// Check for executables count: we should have at least one process (running or exited)
		resp := testUnmarshal[StatusResponse](t, w.Body)
		if resp.Commands.Exited+resp.Commands.Running != 1 {
			t.Fatalf("inconsistency in exited&running balance: %+v", resp.Commands)
		}
	})
}
