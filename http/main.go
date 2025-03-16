package main

import (
	"flag"
	"github.com/brickingsoft/rio"
	"net/http"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "server port")
	flag.Parse()

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
