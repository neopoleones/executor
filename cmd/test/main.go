package main

import (
	"context"
	"executor/internal/executor/naive"
	"executor/internal/models"
	"executor/internal/storage/inmemory"
	"fmt"
	"time"
)

func main() {
	storage, _ := inmemory.GetStorage()
	exec := naive.GetExecutor(storage)

	nr, err := storage.AddCommand(nil, []string{
		"id\n",
		"uname -r\n",
	})

	go func() {
		time.Sleep(time.Second * 4)
		nr.UpdateStatus(models.StatusRejected)

	}()

	if err != nil {
		panic(err)
	}

	res, err := exec.Run(
		context.Background(),
		nr.Sid,
	)
	fmt.Printf("%+v %v\n", res, err)
}
