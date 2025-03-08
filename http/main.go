package main

import (
	"github.com/brickingsoft/rio"
	"net/http"
)

func main() {
	ln, err := rio.Listen("tcp", ":8080")
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
