package internal

import (
	"container/heap"
	"fmt"
	"net"
	"os"
	"time"
)

// PingResult represents single result of TCPing.
type PingResult struct {
	IP          net.IP
	Latency     time.Duration
	Transmitted int
	Received    int
}

type FinalResult []PingResult

func (r FinalResult) Len() int {
	return len(r)
}

func (r FinalResult) Less(i, j int) bool {
	return r[i].Latency < r[j].Latency
}

func (r FinalResult) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (r *FinalResult) Push(x interface{}) {
	*r = append(*r, x.(PingResult))
}

func (r *FinalResult) Pop() interface{} {
	result := *r
	n := len(result)
	item := result[n-1]
	*r = result[0 : n-1]
	return item
}

func (r *FinalResult) Top(top int) {
	var c int
	total := len(*r)
	if total > top {
		c = top
	} else {
		c = total
	}

	fmt.Fprintf(os.Stderr, "\u001B[2K\nTop %d / %d:\n", top, total)
	fmt.Printf("\t%-15s\t%s\t%s\t%s\n", "IP Address", "Lat", "Tran", "Rece")
	for i := 0; i < c; i++ {
		r := heap.Pop(r).(PingResult)
		fmt.Printf("\t%-15s\t%dms\t%d\t%d\n", r.IP, r.Latency.Milliseconds(), r.Transmitted, r.Received)
	}
}
