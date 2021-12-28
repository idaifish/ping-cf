package internal

import (
	"net"
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
