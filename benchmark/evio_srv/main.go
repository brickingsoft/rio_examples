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
		Total data sent:     177.4 MiB (185991168 bytes)
		Total data received: 176.0 MiB (184593536 bytes)
		Bandwidth per channel: 5.925⇅ Mbps (740.6 kBps)
		Aggregate bandwidth: 147.555↓, 148.673↑ Mbps
		Packet rate estimate: 18568.5↓, 12776.1↑ (3↓, 44↑ TCP MSS/op)
		Test duration: 10.0081 s.
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
