package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {

	/*
		Destination: [192.168.100.120]:9000
		Interface eth0 address [192.168.100.1]:0
		Using interface eth0 to connect to [192.168.100.120]:9000
		Ramped up to 50 connections.
		Total data sent:     198.3 MiB (207945728 bytes)
		Total data received: 196.6 MiB (206165284 bytes)
		Bandwidth per channel: 6.623⇅ Mbps (827.9 kBps)
		Aggregate bandwidth: 164.859↓, 166.282↑ Mbps
		Packet rate estimate: 14937.1↓, 14506.0↑ (2↓, 45↑ TCP MSS/op)
		Test duration: 10.0045 s.
	*/
	var port int
	flag.IntVar(&port, "port", 9000, "server port")
	flag.Parse()

	ln, lnErr := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
