package main

import (
	"context"
	"executor/internal/config"
	"executor/internal/storage/postgres"
	"fmt"
)

func main() {
	// test
	cfg := config.GetConfiguration()

	appCtx := context.Background()
	sampleScript := []string{
		"id",
		"uname -a",
	}

	storage, err := postgres.GetStorage(appCtx, cfg)
	if err != nil {
		panic(err)
	}

	if m, e := storage.AddCommand(appCtx, sampleScript); e != nil {
		panic(e)
	} else {
		fmt.Println(m)
	}

	m, e := storage.GetCommands(appCtx)
	if e != nil {
		panic(e)
	}

	for _, cr := range m {
		fmt.Println(cr)
	}

	//storage, _ := inmemory.GetStorage()
	//exec := naive.GetExecutor(storage)
	//
	//nr, err := storage.AddCommand(nil, []string{
	//	"id\n",
	//	"uname -r\n",
	//})
	//
	//go func() {
	//	time.Sleep(time.Second * 4)
	//	nr.UpdateStatus(models.StatusRejected)
	//
	//}()
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//res, err := exec.Run(
	//	context.Background(),
	//	nr.Sid,
	//)
	//fmt.Printf("%+v %v\n", res, err)
}
