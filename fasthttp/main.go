package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/liburing/aio"
	"github.com/valyala/fasthttp"
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
	srvErr := fasthttp.Serve(ln, func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.SetStatusCode(200)
		_, _ = ctx.WriteString("hello world")
	})
	if srvErr != nil {
		panic(srvErr)
		return
	}
}
