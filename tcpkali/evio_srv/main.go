package main

import (
	"flag"
	"fmt"
	"github.com/tidwall/evio"
	"log"
)

func main() {
	/*
		Destination: [192.168.100.120]:9000
		Interface eth0 address [192.168.100.1]:0
		Using interface eth0 to connect to [192.168.100.120]:9000
		Ramped up to 50 connections.
		Total data sent:     200.4 MiB (210108416 bytes)
		Total data received: 198.6 MiB (208234360 bytes)
		Bandwidth per channel: 6.688⇅ Mbps (836.0 kBps)
		Aggregate bandwidth: 166.458↓, 167.956↑ Mbps
		Packet rate estimate: 14272.9↓, 14412.0↑ (2↓, 44↑ TCP MSS/op)
		Test duration: 10.0078 s.
	*/
	var port int
	var loops int
	var reuseport bool

	flag.IntVar(&port, "port", 9000, "server port")
	flag.BoolVar(&reuseport, "reuseport", false, "reuseport (SO_REUSEPORT)")
	flag.IntVar(&loops, "loops", 0, "num loops")
	flag.Parse()

	var events evio.Events
	events.NumLoops = loops
	events.Serving = func(srv evio.Server) (action evio.Action) {
		return
	}
	events.Data = func(c evio.Conn, in []byte) (out []byte, action evio.Action) {
		out = in
		return
	}
	scheme := "tcp"
	log.Fatal(evio.Serve(events, fmt.Sprintf("%s://:%d?reuseport=%t", scheme, port, reuseport)))
}
