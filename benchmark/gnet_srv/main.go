package main

import (
	"flag"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"log"
)

func main() {

	/* tcpkali --workers 1 -c 50 -T 10s -m "PING" 192.168.100.120:9000
	Destination: [192.168.100.120]:9000
	Interface eth0 address [192.168.100.1]:0
	Using interface eth0 to connect to [192.168.100.120]:9000
	Ramped up to 50 connections.
	Total data sent:     183.8 MiB (192741376 bytes)
	Total data received: 182.3 MiB (191161224 bytes)
	Bandwidth per channel: 6.136⇅ Mbps (767.0 kBps)
	Aggregate bandwidth: 152.776↓, 154.039↑ Mbps
	Packet rate estimate: 18598.8↓, 13340.6↑ (3↓, 44↑ TCP MSS/op)
	Test duration: 10.01 s.
	*/

	/* tcpkali --workers 1 -c 50 -r 5k -m "PING" 192.168.100.120:9000
	Destination: [192.168.100.120]:9000
	Interface eth0 address [192.168.100.1]:0
	Using interface eth0 to connect to [192.168.100.120]:9000
	Ramped up to 50 connections.
	Total data sent:     9.5 MiB (10011192 bytes)
	Total data received: 9.5 MiB (10011192 bytes)
	Bandwidth per channel: 0.320⇅ Mbps (40.0 kBps)
	Aggregate bandwidth: 8.008↓, 8.008↑ Mbps
	Packet rate estimate: 28936.6↓, 28957.4↑ (1↓, 1↑ TCP MSS/op)
	Test duration: 10.0007 s.
	*/

	var port int
	var multicore bool

	flag.IntVar(&port, "port", 9000, "--port 9000")
	flag.BoolVar(&multicore, "multicore", false, "--multicore true")
	flag.Parse()
	echo := &echoServer{addr: fmt.Sprintf("tcp://:%d", port), multicore: multicore}
	log.Fatal(gnet.Run(echo, echo.addr, gnet.WithMulticore(multicore)))

}

type echoServer struct {
	gnet.BuiltinEventEngine

	eng       gnet.Engine
	addr      string
	multicore bool
}

func (es *echoServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	return gnet.None
}

func (es *echoServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
}
