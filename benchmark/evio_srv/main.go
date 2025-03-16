package main

import (
	"flag"
	"fmt"
	"github.com/tidwall/evio"
	"log"
)

func main() {
	/* tcpkali --workers 1 -c 50 -T 10s -m "PING" 192.168.100.120:9000
	Destination: [192.168.100.120]:9000
	Interface eth0 address [192.168.100.1]:0
	Using interface eth0 to connect to [192.168.100.120]:9000
	Ramped up to 50 connections.
	Total data sent:     182.6 MiB (191496192 bytes)
	Total data received: 181.1 MiB (189878896 bytes)
	Bandwidth per channel: 6.100⇅ Mbps (762.5 kBps)
	Aggregate bandwidth: 151.862↓, 153.156↑ Mbps
	Packet rate estimate: 19010.4↓, 13192.0↑ (3↓, 44↑ TCP MSS/op)
	Test duration: 10.0027 s.
	*/

	/* tcpkali --workers 1 -c 50 -r 5k -m "PING" 192.168.100.120:9000
	Destination: [192.168.100.120]:9000
	Interface eth0 address [192.168.100.1]:0
	Using interface eth0 to connect to [192.168.100.120]:9000
	Ramped up to 50 connections.
	Total data sent:     9.5 MiB (10011756 bytes)
	Total data received: 9.5 MiB (10011756 bytes)
	Bandwidth per channel: 0.320⇅ Mbps (40.0 kBps)
	Aggregate bandwidth: 8.009↓, 8.009↑ Mbps
	Packet rate estimate: 29327.7↓, 29375.9↑ (1↓, 1↑ TCP MSS/op)
	Test duration: 10.0011 s.
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
