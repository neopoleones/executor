package main

import (
	"context"
	"executor/internal/config"
	"executor/internal/executor/naive"
	"executor/internal/service/rest"
	"executor/internal/storage"
	"executor/internal/storage/inmemory"
	"executor/internal/storage/postgres"
	"os"
	"os/signal"
	"syscall"
)

var exitSignals = []os.Signal{
	syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
}

func main() {
	var err error
	var st storage.ExecutorStorage

	appCtx, _ := signal.NotifyContext(context.Background(), exitSignals...)
	cfg := config.GetConfiguration()

	// Get storage
	switch cfg.Database.Kind {
	case config.DBKindLocal:
		st, err = inmemory.GetStorage()
		if err != nil {
			panic(err)
		}
	case config.DBKindPostgres:
		st, err = postgres.GetStorage(appCtx, cfg)
		if err != nil {
			panic(err)
		}
	default:
		panic("not implemented")
	}

	// Get executor
	runner := naive.GetExecutor(st, cfg)

	// Create & setup service
	service := rest.GetService(cfg)
	service.Setup(st)

	// Use context and run
	go runner.Start(appCtx)
	defer runner.Release(appCtx)

	if err := service.Run(appCtx); err != nil {
		panic(err)
	}
}
