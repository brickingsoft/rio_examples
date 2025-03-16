package main

import (
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio_examples/benchmark/local_bench_tcp"
	"github.com/brickingsoft/rio_examples/images"
	"github.com/panjf2000/gnet/v2"
	"github.com/tidwall/evio"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	/* print
	------net------
	Total data sent: 1.7M (1807572 bytes)
	Total data received: 1.7M (1807572 bytes)
	sent/sec: 180707.27
	recv/sec: 180707.27

	------gnet------
	Total data sent: 1.8M (1883164 bytes)
	Total data received: 1.8M (1883164 bytes)
	sent/sec: 188224.41
	recv/sec: 188224.41

	------evio------
	Total data sent: 1.7M (1797752 bytes)
	Total data received: 1.7M (1797752 bytes)
	sent/sec: 179709.45
	recv/sec: 179709.45

	------rio------
	Total data sent: 6.6M (6950212 bytes)
	Total data received: 6.6M (6950212 bytes)
	sent/sec: 694868.01
	recv/sec: 694868.01
	*/

	rio.Pin()
	defer rio.Unpin()

	var (
		values = make([]float64, 0, 1)
		names  = make([]string, 0, 1)
		out    = "benchmark/out/bench_local.png"
	)

	srvs := []local_bench_tcp.Serve{
		serveNet,
		serveGnet,
		serveEvio,
		serveRIO,
	}
	port := 9000
	for _, srv := range srvs {
		port++
		rm, rmErr := local_bench_tcp.Bench(port, 50, 10*time.Second, srv)
		if rmErr != nil {
			fmt.Println("err:", rmErr)
			return
		}
		fmt.Println(rm)
		values = append(values, rm.Rate())
		names = append(names, rm.Title())
	}

	images.Plotit(out, "Echo(C50 T10s)", values, names)
}

func serveRIO(port int) (title string, closer io.Closer, err error) {
	title = "rio"
	ln, lnErr := rio.Listen("tcp", fmt.Sprintf(":%d", port))
	if lnErr != nil {
		err = lnErr
		return
	}
	closer = ln
	go func(ln net.Listener) {
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
				return
			}
			go func(conn net.Conn) {
				var packet [0xFFF]byte
				for {
					rn, rErr := conn.Read(packet[:])
					if rErr != nil {
						conn.Close()
						return
					}
					_, wEr := conn.Write(packet[:rn])
					if wEr != nil {
						conn.Close()
						return
					}
				}
			}(conn)
		}
	}(ln)
	return
}

func serveNet(port int) (title string, closer io.Closer, err error) {
	title = "net"
	ln, lnErr := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if lnErr != nil {
		err = lnErr
		return
	}
	closer = ln
	go func(ln net.Listener) {
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
				return
			}
			go func(conn net.Conn) {
				var packet [0xFFF]byte
				for {
					conn.SetDeadline(time.Now().Add(15 * time.Second))
					rn, rErr := conn.Read(packet[:])
					if rErr != nil {
						conn.Close()
						return
					}
					_, wEr := conn.Write(packet[:rn])
					if wEr != nil {
						conn.Close()
						return
					}
				}
			}(conn)
		}
	}(ln)
	return
}

type emptyCloser struct{}

func (e emptyCloser) Close() error {
	return nil
}

func serveEvio(port int) (title string, closer io.Closer, err error) {
	go func(port int) {
		var events evio.Events
		events.NumLoops = 1
		events.Serving = func(srv evio.Server) (action evio.Action) {
			return
		}
		events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
			out = in
			return
		}

		scheme := "tcp"
		log.Fatal(evio.Serve(events, fmt.Sprintf("%s://:%d?reuseport=%t", scheme, port, false)))
	}(port)
	title = "evio"
	closer = emptyCloser{}
	time.Sleep(50 * time.Millisecond)
	return
}

func serveGnet(port int) (title string, closer io.Closer, err error) {
	go func(port int) {
		echo := &gnetServer{addr: fmt.Sprintf("tcp://:%d", port), multicore: true}
		log.Fatal(gnet.Run(echo, echo.addr, gnet.WithMulticore(true), gnet.WithLogger(&gnetLogger{})))
	}(port)
	title = "gnet"
	closer = emptyCloser{}
	time.Sleep(50 * time.Millisecond)
	return
}

type gnetServer struct {
	gnet.BuiltinEventEngine
	eng       gnet.Engine
	addr      string
	multicore bool
}

func (es *gnetServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	return gnet.None
}

func (es *gnetServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
}

type gnetLogger struct{}

func (g *gnetLogger) Debugf(format string, args ...any) {
	return
}

func (g *gnetLogger) Infof(format string, args ...any) {
	return
}

func (g *gnetLogger) Warnf(format string, args ...any) {
	return
}

func (g *gnetLogger) Errorf(format string, args ...any) {
	return
}

func (g *gnetLogger) Fatalf(format string, args ...any) {
	return
}
