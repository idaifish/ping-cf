package main

import (
	"container/heap"
	"context"
	"fmt"
	. "github.com/idaifish/ping-cf/internal"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
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

	for i := 0; i < 64; i++ {
		wg.Add(1)
		go Worker(ctx, wg, cloudflareIPs, resultsChan)
	}

	wg.Wait()

	fmt.Fprintf(os.Stderr, "\u001B[2K\nTop 10 / %d:\n", len(results))

	if len(results) > 10 {
		fmt.Printf("\t%-15s\t%s\t%s\t%s\n", "IP Address", "Lat", "Tran", "Rece")
		for i := 0; i < 10; i++ {
			r := heap.Pop(&results).(PingResult)
			fmt.Printf("\t%-15s\t%dms\t%d\t%d\n", r.IP, r.Latency.Milliseconds(), r.Transmitted, r.Received)
		}
	}
}
