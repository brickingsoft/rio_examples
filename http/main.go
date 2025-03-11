package main

import (
	"flag"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/iouring/aio"
	"net/http"
)

func main() {
	var port int
	var schema string
	flag.IntVar(&port, "port", 9000, "server port")
	flag.StringVar(&schema, "schema", aio.DefaultFlagsSchema, "iouring schema")
	flag.Parse()

	rio.PrepareIOURingSetupOptions(
		aio.WithFlagsSchema(schema),
	)
	ln, err := rio.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
		return
	}
	defer ln.Close()
	http.Serve(ln, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("hello world"))
	}))
}
