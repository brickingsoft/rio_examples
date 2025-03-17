package local

import (
	"fmt"
	"github.com/brickingsoft/rio"
	"io"
	"net"
	"sync"
	"time"
)

type Serve func(port int) (title string, closer io.Closer, err error)

type Dialer func(network string, address string) (conn net.Conn, err error)

var (
	RioDialer Dialer = rio.Dial
	NetDialer Dialer = net.Dial
)

func Bench(port int, count int, dur time.Duration, srvFn Serve, dialer Dialer) (met *Metric, err error) {
	title, srv, srvErr := srvFn(port)
	if srvErr != nil {
		return nil, srvErr
	}
	defer srv.Close()

	met = NewMetric(title)
	wg := new(sync.WaitGroup)
	wg.Add(count)
	met.Begin()
	for i := 0; i < count; i++ {
		go func(wg *sync.WaitGroup, port int, dur time.Duration, met *Metric) {
			defer wg.Done()
			conn, connErr := dialer("tcp", fmt.Sprintf(":%d", port))
			if connErr != nil {
				fmt.Println("conn err:", connErr)
				return
			}
			timer := time.NewTimer(dur)
			stopped := false
			sp := []byte("PING")
			rp := make([]byte, 1024)
			for {
				select {
				case <-timer.C:
					stopped = true
					break
				default:
					wn, wEr := conn.Write(sp)
					if wEr != nil {
						conn.Close()
						return
					}
					met.IncrOut(wn)
					rn, rErr := conn.Read(rp)
					if rErr != nil {
						conn.Close()
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
		}(wg, port, dur, met)
	}
	wg.Wait()
	met.End()
	return
}
