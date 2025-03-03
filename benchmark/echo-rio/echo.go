package echorio

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/brickingsoft/rio"
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

	setting()

	ln, lnErr := serve(port, nBytes)
	if lnErr != nil {
		err = lnErr
		return
	}
	time.Sleep(1 * time.Second)
	met := metric.New()
	met.Begin()
	dial(met, workers, count, port, nBytes)
	met.End()
	_ = ln.Close()
	actions, inbounds, outbounds = met.PerSecond()
	failures = met.Failures()
	dur = met.CostDuration()
	return
}

func serve(port int, nBytes int) (ln net.Listener, err error) {
	ln, err = rio.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return
	}
	go func(ln net.Listener, nBytes int) {
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
				break
			}
			go func(conn net.Conn) {
				b := make([]byte, nBytes)
				for {
					rn, rErr := conn.Read(b)
					if rErr != nil {
						_ = conn.Close()
						if errors.Is(rErr, io.EOF) {
							break
						}
						break
					}

					_, wErr := conn.Write(b[:rn])
					if wErr != nil {
						_ = conn.Close()
						break
					}
				}
			}(conn)
		}
	}(ln, nBytes)
	return
}

func dial(met *metric.Metric, workers int, count int, port int, nBytes int) {
	b := make([]byte, nBytes)
	_, _ = rand.Read(b)
	wg := new(sync.WaitGroup)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(met *metric.Metric, wg *sync.WaitGroup, count int, port int, b []byte) {
			defer wg.Done()
			for j := 0; j < count; j++ {
				conn, err := net.Dial("tcp", fmt.Sprintf("127.0.0.1:%d", port))
				if err != nil {
					time.Sleep(10 * time.Microsecond)
					j--
					continue
				}
				met.IncACT(1)
				remain := len(b)
				for remain > 0 {
					wn, wErr := conn.Write(b)
					met.IncOUT(wn)
					if wErr != nil {
						_ = conn.Close()
						met.Failed(1)
						return
					}
					remain -= wn
				}
				rn, rErr := conn.Read(b)
				met.IncIN(rn)
				if rErr != nil {
					_ = conn.Close()
					met.Failed(1)
					return
				}
				_ = conn.Close()
			}
		}(met, wg, count, port, b)
	}
	wg.Wait()
}
