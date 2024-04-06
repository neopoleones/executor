package rest

import (
	"context"
	"executor/internal/config"
	"executor/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type Service struct {
	addr    string
	version string

	storage storage.ExecutorStorage

	server http.Server
}

func (s *Service) Release(ctx context.Context) {
	err := s.server.Shutdown(ctx)

	// No reason to return error: graceful shutdown failed!
	if err != nil {
		panic(err)
	}

	s.storage.Close(ctx)
}

func (s *Service) Run(ctx context.Context) error {

	// I guess I can use background here as after closing ctx.Done we'll exit or call serverStopCtx
	// So serverCtx(background) is anyway related to ctx(outbound background)

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

		s.Release(shutdownCtx)
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
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)
}

func (s *Service) setupHandlers(r *chi.Mux) {
	r.Get("/status", s.statusHandler())

	apiRouter := chi.NewRouter()

	apiRouter.Get("/get", s.getHandler())
	apiRouter.Get("/list", s.listHandler())
	apiRouter.Post("/schedule", s.scheduleHandler())
	apiRouter.Post("/kill", s.killHandler())

	r.Mount("/cmd", apiRouter)
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

func GetService(cfg *config.Configuration) *Service {
	return &Service{
		version: cfg.Version,
		addr:    cfg.Service.Addr,
	}
}
