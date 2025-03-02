package echostd

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/brickingsoft/rio_examples/benchmark/metric"
	"io"
	"net"
	"sync"
	"time"
)

func Bench(workers int, count int, port int, nBytes int) (dur time.Duration, actions uint64, inbounds uint64, outbounds uint64, failures uint64, err error) {
	if workers < 1 {
		workers = 1
	}
	if count < 1 {
		count = 1
	}
	if port < 1 {
		port = 9000
	}
	if nBytes < 1 {
		nBytes = 1024
	}

	met := metric.New()

	ln, lnErr := serve(met, port, nBytes)
	if lnErr != nil {
		err = lnErr
		return
	}
	time.Sleep(1 * time.Second)
	met.Begin()
	dial(workers, count, port, nBytes)
	met.End()
	_ = ln.Close()
	actions, inbounds, outbounds = met.PerSecond()
	failures = met.Failures()
	dur = met.CostDuration()
	return
}

func serve(met *metric.Metric, port int, nBytes int) (ln net.Listener, err error) {
	ln, err = net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	go func(ln net.Listener, met *metric.Metric, nBytes int) {
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
				return
			}
			met.IncACT(1)
			go func(conn net.Conn, met *metric.Metric) {
				b := make([]byte, nBytes)
				for {
					rn, rErr := conn.Read(b)
					if rErr != nil {
						_ = conn.Close()
						if errors.Is(rErr, io.EOF) {
							return
						}
						met.Failed(1)
						return
					}
					met.IncIN(rn)

					wn, wErr := conn.Write(b[:rn])
					if wErr != nil {
						_ = conn.Close()
						met.Failed(1)
						return
					}
					met.IncOUT(wn)
				}
			}(conn, met)
		}
	}(ln, met, nBytes)
	return
}

func dial(workers int, count int, port int, nBytes int) {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	wg := new(sync.WaitGroup)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(wg *sync.WaitGroup, count int, port int, b []byte) {
			defer wg.Done()
			for j := 0; j < count; j++ {
				conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
				if err != nil {
					time.Sleep(10 * time.Microsecond)
					continue
				}
				remain := len(b)
				for remain > 0 {
					wn, wErr := conn.Write(b)
					if wErr != nil {
						_ = conn.Close()
						return
					}
					remain -= wn
				}
				_, _ = conn.Read(b)
				_ = conn.Close()
			}
		}(wg, count, port, b)
	}
	wg.Wait()
}
