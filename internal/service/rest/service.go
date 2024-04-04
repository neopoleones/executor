package rest

import (
	"context"
	"executor/internal/storage"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Service struct {
	addr    string
	storage storage.ExecutorStorage
	server  http.Server
}

func (s *Service) Run(ctx context.Context) error {
	serverCtx, serverStopCtx := context.WithCancel(context.TODO())

	go func() {
		// Exit signal sent (or any other reason); We need to prepare server for that
		<-ctx.Done()
		slog.Info("Exiting gracefully")

		shutdownCtx, shutdownOk := context.WithTimeout(serverCtx, 20*time.Second)
		defer shutdownOk()

		go func() {
			<-shutdownCtx.Done()

			if shutdownCtx.Err() == context.DeadlineExceeded {
				// Server failed to shutdown in 20 seconds. Force killing
				slog.Error("Server failed to shutdown. Force exit")
				os.Exit(-1)
			}

		}()

		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			// No reason to return error: graceful shutdown failed!
			panic(err)
		}

		serverStopCtx()
	}()

	slog.Info(
		"Starting the REST service",
		slog.String("address", s.addr),
	)

	err := s.server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}

	// Wait for clean-up
	<-serverCtx.Done()
	return nil
}

func (s *Service) setupMiddlewares(r *chi.Mux) {
	// TODO
}

func (s *Service) setupHandlers(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("Hello from teapot"))
	})
}

func (s *Service) getRouter() http.Handler {
	router := chi.NewRouter()

	s.setupMiddlewares(router)
	s.setupHandlers(router)

	return router
}

func (s *Service) Setup(es storage.ExecutorStorage) {
	s.storage = es

	s.server = http.Server{
		Addr:    s.addr,
		Handler: s.getRouter(),
	}
}

func GetService(addr string) *Service {
	return &Service{
		addr: addr,
	}
}
