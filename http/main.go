package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/liburing/aio"
	"net/http"
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
	srvErr := http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("hello world"))
	}))
	if srvErr != nil {
		panic(srvErr)
	}
}
