package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"log"
	"net"
)

func main() {
	/*
		Destination: [192.168.100.120]:9000
		Interface eth0 address [192.168.100.1]:0
		Using interface eth0 to connect to [192.168.100.120]:9000
		Ramped up to 50 connections.
		Total data sent:     287.6 MiB (301548988 bytes)
		Total data received: 286.4 MiB (300361173 bytes)
		Bandwidth per channel: 9.627⇅ Mbps (1203.3 kBps)
		Aggregate bandwidth: 240.188↓, 241.138↑ Mbps
		Packet rate estimate: 27791.8↓, 20820.1↑ (3↓, 32↑ TCP MSS/op)
		Test duration: 10.0042 s.
	*/

	var port int
	flag.IntVar(&port, "port", 9000, "server port")
	flag.Parse()

	ln, lnErr := rio.Listen("tcp", fmt.Sprintf(":%d", port))
	if lnErr != nil {
		log.Fatal("lnErr:", lnErr)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go func(conn net.Conn) {
			var packet [0xFFF]byte
			for {
				rn, rErr := conn.Read(packet[:])
				if rErr != nil {
					conn.Close()
					return
				}
				_, wEr := conn.Write(packet[:rn])
				if wEr != nil {
					conn.Close()
					return
				}
			}
		}(conn)
	}
}
