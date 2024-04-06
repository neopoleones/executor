package main

import (
	"context"
	"executor/internal/config"
	"executor/internal/executor/naive"
	"executor/internal/service/rest"
	"executor/internal/storage"
	"executor/internal/storage/inmemory"
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

	cfg := config.GetConfiguration()

	// Get storage
	switch cfg.Database.Kind {
	case config.DBKindLocal:
		st, err = inmemory.GetStorage()
		if err != nil {
			panic(err)
		}
	case config.DBKindPostgres:
		fallthrough
	default:
		panic("not implemented")
	}

	// Get executor
	runner := naive.GetExecutor(st, cfg)

	// Create & setup service
	service := rest.GetService(cfg)
	service.Setup(st, runner)

	// Prepare context and run
	appCtx, _ := signal.NotifyContext(context.Background(), exitSignals...)
	if err := service.Run(appCtx); err != nil {
		panic(err)
	}
}
