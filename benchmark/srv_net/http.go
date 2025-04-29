package srv_net

import (
	"fmt"
	"net"
	"net/http"
)

func ServeHttp(port int) {
	ln, lnErr := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if lnErr != nil {
		panic(lnErr)
		return
	}

	srv := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("hello world"))
		}),
	}

	go func(srv *http.Server, ln net.Listener) {
		if srvErr := srv.Serve(ln); srvErr != nil {
			panic(srvErr)
		}
	}(srv, ln)
}
