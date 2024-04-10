package rest

import (
	"context"
	"executor/internal/models"
	"executor/internal/storage"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
)

const getEndpoint = "/cmd/get?sid=%s"

type getRequestArgs struct {
	sid string
}

func getRequest(h http.Handler, rd getRequestArgs) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(getEndpoint, rd.sid), nil)
	if err != nil {
		return nil, err
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)

	return w, nil
}

func TestGetHandlerIncorrectCommand(t *testing.T) {
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		rn, _ := es.AddCommand(context.Background(), []string{
			"id",
		})

		w := testPerfromRequestWithData(t, h, getRequest, getRequestArgs{sid: rn.Sid.String()})
		assertCode(t, w, http.StatusOK)

		resp := testUnmarshal[GetResponse](t, w.Body)
		assertStatus(t, resp.Status, "ok")

		if resp.Runnable.Sid.String() != rn.Sid.String() {
			t.Fatalf("Incorrect runnable returned (sid)")
		}

		if resp.Runnable.Status != models.StatusScheduled {
			t.Fatalf("Incorrect runnable returned (status)")
		}
	})
}

func TestGetHandlerIncorrectID(t *testing.T) {
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		w := testPerfromRequestWithData(t, h, getRequest, getRequestArgs{sid: uuid.New().String()})
		assertCode(t, w, http.StatusNotFound)
	})
}

func TestGetHandlerBadRequest(t *testing.T) {
	performTest(func(es storage.ExecutorStorage, h http.Handler) {
		w := testPerfromRequestWithData(t, h, getRequest, getRequestArgs{sid: "impossible_to_parse"})
		assertCode(t, w, http.StatusBadRequest)
	})
}
