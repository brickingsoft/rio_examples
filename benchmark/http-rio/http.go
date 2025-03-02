package http_rio

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio_examples/benchmark/metric"
	"io"
	"net"
	"net/http"
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
	ln, err = rio.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return
	}

	go func(ln net.Listener, met *metric.Metric, nBytes int) {
		_ = http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			met.IncACT(1)
			b := make([]byte, nBytes)
			rn, rErr := r.Body.Read(b)
			defer r.Body.Close()
			met.IncIN(rn)
			for rn > 0 {
				wn, wErr := w.Write(b)
				if wErr != nil {
					met.Failed(1)
					break
				}
				met.IncOUT(wn)
				rn -= wn
				b = b[wn:]
			}
			if rErr != nil {
				if rErr == io.EOF {
					return
				}
				met.Failed(1)
				return
			}

		}))
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
				resp, postErr := http.Post(fmt.Sprintf("http://127.0.0.1:%d", port), "application/text", bytes.NewBuffer(b))
				if postErr != nil {
					time.Sleep(10 * time.Microsecond)
					j--
					continue
				}
				_ = resp.Body.Close()
			}
		}(wg, count, port, b)
	}
	wg.Wait()
}
