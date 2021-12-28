package main

import (
	"container/heap"
	"context"
	"flag"
	"fmt"
	. "github.com/idaifish/ping-cf/internal"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var numWorker = flag.Int("n", 64, "number of workers")
var top = flag.Int("top", 10, "print topN results")

func main() {
	flag.Parse()

	var results FinalResult
	heap.Init(&results)
	resultsChan := make(chan PingResult, 100)
	cloudflareIPs := GetIP()

	wg := new(sync.WaitGroup)
	ctx, cancel := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT)
	go func() {
		<-sig
		fmt.Println("\rCancelled by user. Wait. . .")
		cancel()
	}()

	go func() {
		for result := range resultsChan {
			heap.Push(&results, result)
		}
	}()

	for i := 0; i < *numWorker; i++ {
		wg.Add(1)
		go Worker(ctx, wg, cloudflareIPs, resultsChan)
	}

	wg.Wait()

	results.Top(*top)
}
