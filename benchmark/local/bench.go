package local

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio_examples/benchmark/images"
	"github.com/brickingsoft/rio_examples/benchmark/srv_evio"
	"github.com/brickingsoft/rio_examples/benchmark/srv_gnet"
	"github.com/brickingsoft/rio_examples/benchmark/srv_net"
	"github.com/brickingsoft/rio_examples/benchmark/srv_rio"
	"io"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var rioCase = Cause{
	Name:  "RIO",
	Serve: srv_rio.Serve,
	Dial:  rio.Dial,
}

var netCase = Cause{
	Name:  "NET",
	Serve: srv_net.Serve,
	Dial:  rio.Dial,
}

var gnetCase = Cause{
	Name:  "GNET",
	Serve: srv_gnet.Serve,
	Dial:  rio.Dial,
}

var evioCase = Cause{
	Name:  "EVIO",
	Serve: srv_evio.Serve,
	Dial:  rio.Dial,
}

func Bench(port int, count int, dur time.Duration, msg string, out string) {
	cases := []Cause{rioCase, evioCase, gnetCase, netCase}

	ms := make([]*Metric, 0)
	buf := bytes.NewBufferString("")
	for _, c := range cases {
		port++
		met, bErr := c.bench(port, count, dur, msg)
		if bErr != nil {
			fmt.Println(c.Name, "failed")
			fmt.Println(bErr)
		} else {
			fmt.Println(met.String())
			buf.WriteString(fmt.Sprintf("%s\n", met.String()))
			ms = append(ms, met)
		}
	}
	_ = os.MkdirAll(out, 0777)
	// write text
	textFile := filepath.Join(out, "benchmark_local.txt")
	_ = os.WriteFile(textFile, buf.Bytes(), 0777)
	// write image
	req := images.Request{
		Path:  filepath.Join(out, "benchmark_local.png"),
		Title: "Benchmark(LOCAL)",
		Label: "req/s",
		Items: make([]images.Item, 0, 1),
	}
	for _, m := range ms {
		req.Items = append(req.Items, images.Item{
			Name:  m.Title(),
			Value: m.Rate(),
		})
	}

	err := images.Draw(req)
	if err != nil {
		fmt.Println(err)
	}
	return
}

type Serve func(port int)

type Dialer func(network string, address string) (conn net.Conn, err error)

type Cause struct {
	Name  string
	Serve Serve
	Dial  Dialer
}

func (c Cause) bench(port int, count int, dur time.Duration, msg string) (met *Metric, err error) {
	c.Serve(port)
	time.Sleep(50 * time.Millisecond)

	host := "127.0.0.1"
	met = NewMetric(c.Name, len(msg))
	wg := new(sync.WaitGroup)
	wg.Add(count)
	met.Begin()
	for i := 0; i < count; i++ {
		go func(wg *sync.WaitGroup, host string, port int, dur time.Duration, msg string, met *Metric) {
			defer wg.Done()
			conn, connErr := c.Dial("tcp", fmt.Sprintf("%s:%d", host, port))
			if connErr != nil {
				met.IncrFailed()
				return
			}
			_ = conn.SetDeadline(time.Now().Add(dur))
			timer := time.NewTimer(dur)
			stopped := false
			sp := []byte(msg)
			rp := make([]byte, 1024)
			for {
				select {
				case <-timer.C:
					stopped = true
					break
				default:
					wn, wEr := conn.Write(sp)
					if wEr != nil {
						if !errors.Is(wEr, net.ErrClosed) && !errors.Is(wEr, context.DeadlineExceeded) {
							met.IncrFailed()
						}
						_ = conn.Close()
						return
					}
					met.IncrOut(wn)
					rn, rErr := conn.Read(rp)
					if rErr != nil {
						if !errors.Is(rErr, io.EOF) && !errors.Is(rErr, context.DeadlineExceeded) {
							met.IncrFailed()
						}
						_ = conn.Close()
						return
					}
					met.IncrIn(rn)
					break
				}
				if stopped {
					break
				}
			}
			_ = conn.Close()
		}(wg, host, port, dur, msg, met)
	}
	wg.Wait()
	met.End()
	return
}
