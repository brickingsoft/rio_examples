package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/valyala/fasthttp"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "server port")
	flag.Parse()

	ln, err := rio.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
		return
	}
	defer ln.Close()
	fasthttp.Serve(ln, func(ctx *fasthttp.RequestCtx) {
		ctx.SetContentType("text/html; charset=utf-8")
		ctx.SetStatusCode(200)
		ctx.WriteString("hello world")
	})
}
