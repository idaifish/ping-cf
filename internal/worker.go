package internal

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

const maxEndurableLatency = 350 * time.Millisecond

var pingCount = flag.Int("c", 4, "count of ping.")
var errUnreachable = errors.New("unreachable IP")

func tcpPing(ip net.IP, count uint8) (result *PingResult, err error) {
	result = &PingResult{
		Latency:     0,
		Transmitted: 0,
		Received:    0,
		IP:          ip,
	}
	fmt.Printf("\033[2K\rPing %s", ip.String())
	var mean time.Duration = 0

	for i := uint8(0); i < count; i++ {
		result.Transmitted += 1
		start := time.Now()
		conn, err := net.DialTimeout("tcp", ip.String()+":443", maxEndurableLatency)
		if err != nil {
			continue
		}
		result.Received += 1
		mean += time.Since(start)
		_ = conn.Close()
	}

	if result.Received > 0 {
		// calculate Latency time.
		result.Latency = mean / time.Duration(result.Received)
		return
	} else {
		return nil, errUnreachable
	}
}

func Worker(ctx context.Context, wg *sync.WaitGroup, ipChan chan net.IP, resultChan chan PingResult) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case ip := <-ipChan:
			if ip != nil {
				result, err := tcpPing(ip, uint8(*pingCount))
				if err == nil {
					if (float32(result.Received) / float32(result.Transmitted)) > 0.75 {
						resultChan <- *result
					}
				}
				continue
			}
			// ipChan is closed.
			return
		}
	}
}
