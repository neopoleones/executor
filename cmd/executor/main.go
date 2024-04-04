package main

import (
	"context"
	"executor/internal/service/rest"
	"executor/internal/storage/inmemory"
	"os"
	"os/signal"
	"syscall"
)

var exitSignals = []os.Signal{
	syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
}

func init() {

}

func main() {
	storage, err := inmemory.GetStorage()
	if err != nil {
		panic(err)
	}

	rs := rest.GetService("127.0.0.1:8080")
	rs.Setup(storage)

	appCtx, _ := signal.NotifyContext(context.Background(), exitSignals...)

	if err := rs.Run(appCtx); err != nil {
		panic(err)
	}
}
