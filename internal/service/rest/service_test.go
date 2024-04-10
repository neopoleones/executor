package rest

import (
	"context"
	"encoding/json"
	"executor/internal/config"
	"executor/internal/executor/naive"
	"executor/internal/storage"
	"executor/internal/storage/inmemory"
	"github.com/google/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func setupForTests() (*Service, *naive.SystemExecutor) {
	// Get configuration and patch for local db
	cfg := &config.Configuration{
		Version: "test",
		Verbose: false,
		Database: config.DatabaseConfiguration{
			Kind: config.DBKindLocal,
		},
		Executor: config.ExecutorConfiguration{
			InterpreterPath: "/bin/sh",
			SchedTicks:      time.Millisecond * 100,
		},
	}
	storage, _ := inmemory.GetStorage()

	exec := naive.GetExecutor(storage, cfg)

	// Initialize service
	srv := GetService(cfg)
	srv.Setup(storage)

	return srv, exec
}

func performTest(performTest func(es storage.ExecutorStorage, h http.Handler)) {
	srv, exec := setupForTests()
	h := srv.server.Handler

	appCtx, appDone := context.WithCancel(context.Background())

	defer func() {
		appDone()
		exec.Release(appCtx)
	}()

	go func() {
		exec.Start(context.Background())
	}()

	performTest(srv.storage, h)
}

type Rec *httptest.ResponseRecorder
type requestFunc func(h http.Handler) (*httptest.ResponseRecorder, error)
type requestData[DT any] func(h http.Handler, data DT) (*httptest.ResponseRecorder, error)

func testPerformRequest(t *testing.T, h http.Handler, requester requestFunc) Rec {
	t.Helper()

	w, err := requester(h)
	if err != nil {
		t.Fatalf("Failed to perfrom request: %v", err)
	}

	return w
}

func testPerfromRequestWithData[DT any](t *testing.T, h http.Handler, requester requestData[DT], data DT) Rec {
	t.Helper()

	w, err := requester(h, data)
	if err != nil {
		t.Fatalf("Failed to perfrom request: %v", err)
	}

	return w
}

func testUnmarshal[RT any](t *testing.T, r io.Reader) RT {
	t.Helper()

	var resp RT
	if err := json.NewDecoder(r).Decode(&resp); err != nil {
		t.Fatalf("bad output format: %v", err)
	}

	return resp
}

func assertCode(t *testing.T, r *httptest.ResponseRecorder, code int) {
	t.Helper()

	if r.Code != code {
		t.Fatalf("bad return code: %v!=%v", r.Code, code)
	}
}

func assertStatus(t *testing.T, s1, s2 string) {
	t.Helper()

	if s1 != s2 {
		t.Fatalf("bad status: %s!=%s", s1, s2)
	}
}

func testParseUUID(t *testing.T, sid string) uuid.UUID {
	t.Helper()

	if v, err := uuid.Parse(sid); err != nil {
		t.Fatalf("failed to parse uuid: %v", err)
		return uuid.Nil
	} else {
		return v
	}
}
