package main

import (
	"context"
	"log"
	"sync"
	"testing"
	"time"
)

func TestMainFunc(t *testing.T) {
	//sigs := make(chan os.Signal, 1)
	//p, err := os.FindProcess(os.Getpid())
	//if err != nil {
	//	panic(err)
	//}

	var wg sync.WaitGroup

	ctx := context.Background()
	cctx, cancelFunc := context.WithCancel(ctx)

	wg.Add(1)
	go func() {
		defer wg.Done()
		run(cctx)
	}()

	// signal is not worked on test due to foreground process receives first
	//signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	//p.Signal(syscall.SIGINT)
	//p.Signal(syscall.SIGTERM)
	//p.Signal(syscall.SIGHUP)

	time.Sleep(2 * time.Second)
	cancelFunc()

	wg.Wait()

	log.Printf("done")
}
