package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/liburing/aio"
	"github.com/valyala/fasthttp"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "server port")
	flag.Parse()
	rio.Preset(
		aio.WithNAPIBusyPollTimeout(time.Microsecond * 50),
	)
	ln, err := rio.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
		return
	}

	srv := &fasthttp.Server{
		Handler: func(ctx *fasthttp.RequestCtx) {
			ctx.SetContentType("text/html; charset=utf-8")
			ctx.SetStatusCode(200)
			_, _ = ctx.WriteString("hello world")
		},
	}
	done := make(chan struct{}, 1)
	go func(ln net.Listener, srv *fasthttp.Server, done chan<- struct{}) {
		if srvErr := srv.Serve(ln); srvErr != nil {
			if errors.Is(srvErr, io.EOF) {
				close(done)
				return
			}
			panic(srvErr)
			return
		}
		close(done)
	}(ln, srv, done)

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	<-signalCh

	if shutdownErr := srv.Shutdown(); shutdownErr != nil {
		panic(shutdownErr)
	}
	<-done
}
