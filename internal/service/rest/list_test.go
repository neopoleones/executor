package rest

import (
	"context"
	"executor/internal/storage"
	"net/http"
	"net/http/httptest"
	"testing"
)

const listEndpoint = "/cmd/list"

func listRequest(h http.Handler) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(http.MethodGet, listEndpoint, nil)
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	return w, nil
}

func TestListHandlerNoCommands(t *testing.T) {
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		w := testPerformRequest(t, h, listRequest)
		assertCode(t, w, http.StatusOK)

		resp := testUnmarshal[ListResponse](t, w.Body)
		assertStatus(t, resp.Status, "ok")

		if len(resp.Commands) != 0 {
			t.Fatalf("storage has commands? len != 0")
		}
	})
}

func TestListHandlerSeveralCommandsExist(t *testing.T) {
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		testCtx := context.Background()

		_, _ = es.AddCommand(testCtx, []string{"id"})

		w := testPerformRequest(t, h, listRequest)
		assertCode(t, w, http.StatusOK)

		resp := testUnmarshal[ListResponse](t, w.Body)
		assertStatus(t, resp.Status, "ok")

		if len(resp.Commands) != 1 {
			t.Fatalf("storage has commands? len != 1")
		}
	})
}
