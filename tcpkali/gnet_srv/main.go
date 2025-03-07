package main

import (
	"flag"
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"log"
)

func main() {
	/*
		Destination: [192.168.100.120]:9000
		Interface eth0 address [192.168.100.1]:0
		Using interface eth0 to connect to [192.168.100.120]:9000
		Ramped up to 50 connections.
		Total data sent:     219.4 MiB (230096896 bytes)
		Total data received: 217.7 MiB (228243396 bytes)
		Bandwidth per channel: 7.329⇅ Mbps (916.1 kBps)
		Aggregate bandwidth: 182.481↓, 183.963↑ Mbps
		Packet rate estimate: 22095.3↓, 15777.4↑ (3↓, 44↑ TCP MSS/op)
		Test duration: 10.0062 s.
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
