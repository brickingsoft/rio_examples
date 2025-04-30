package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/liburing/aio"
	"net"
	"net/http"
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
	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("hello world"))
		}),
	}
	done := make(chan struct{}, 1)
	go func(ln net.Listener, srv *http.Server, done chan<- struct{}) {
		if srvErr := srv.Serve(ln); srvErr != nil {
			if errors.Is(srvErr, http.ErrServerClosed) {
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

	if shutdownErr := srv.Shutdown(context.Background()); shutdownErr != nil {
		panic(shutdownErr)
	}
	<-done
}
