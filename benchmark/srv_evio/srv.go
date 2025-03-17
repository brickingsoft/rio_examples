package srv_evio

import (
	"fmt"
	"github.com/tidwall/evio"
	"log"
)

func Serve(port int) {
	go func(port int) {
		var events evio.Events
		events.NumLoops = 1
		events.Serving = func(srv evio.Server) (action evio.Action) {
			return
		}
		events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
			out = in
			return
		}

		scheme := "tcp"
		log.Fatal(evio.Serve(events, fmt.Sprintf("%s://:%d?reuseport=%t", scheme, port, false)))
	}(port)
	return
}
