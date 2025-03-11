package main

import (
	"github.com/brickingsoft/rio"
	"github.com/valyala/fasthttp"
)

func main() {
	ln, err := rio.Listen("tcp", ":9000")
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
