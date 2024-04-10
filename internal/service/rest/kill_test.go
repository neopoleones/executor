package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"executor/internal/models"
	"executor/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const killEndpoint = "/cmd/kill"

func killRequest(h http.Handler, rd KillRequest) (*httptest.ResponseRecorder, error) {
	reqData, err := json.Marshal(rd)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, killEndpoint, bytes.NewReader(reqData))
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	return w, nil
}

func TestKillHandlerBadCommand(t *testing.T) {
	// Scheduled command can't be killed
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		testCtx := context.Background()
		r, _ := es.AddCommand(testCtx, []string{"id"})

		w := testPerfromRequestWithData[KillRequest](t, h, killRequest, KillRequest{
			RawSid: r.Sid.String(),
		})
		assertCode(t, w, http.StatusUnprocessableEntity)
	})

	// Can't kill command if incorrect id specified
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		w := testPerfromRequestWithData[KillRequest](t, h, killRequest, KillRequest{
			RawSid: "bad_id",
		})
		assertCode(t, w, http.StatusBadRequest)
	})
}

func TestKillHandler(t *testing.T) {
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		t.Parallel()

		appCtx := context.Background()

		w := testPerfromRequestWithData[ScheduleRequest](t, h, scheduleRequest, ScheduleRequest{[]string{
			"while [ 1 ]\ndo\n    echo 'This is an infinite loop' > /dev/null\ndone",
		}})

		resp := testUnmarshal[ScheduleResponse](t, w.Body)
		sid := testParseUUID(t, resp.Sid)

		// Wait for scheduler
		time.Sleep(time.Millisecond * 350)

		// Try to kill
		w = testPerfromRequestWithData(t, h, killRequest, KillRequest{RawSid: sid.String()})
		assertCode(t, w, http.StatusOK)

		killResp := testUnmarshal[KillResponse](t, w.Body)
		assertStatus(t, killResp.Status, "ok")

		// Give time for scheduler again
		time.Sleep(time.Millisecond * 350)

		r, _ := es.GetCommandByID(appCtx, sid)
		if r.Status != models.StatusRejected {
			t.Fatalf("Command not in rejected state")
		}
	})
}
