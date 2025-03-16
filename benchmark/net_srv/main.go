package main

import (
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {

	/* tcpkali --workers 1 -c 50 -T 10s -m "PING" 192.168.100.120:9000
	Destination: [192.168.100.120]:9000
	Interface eth0 address [192.168.100.1]:0
	Using interface eth0 to connect to [192.168.100.120]:9000
	Ramped up to 50 connections.
	Total data sent:     183.4 MiB (192282624 bytes)
	Total data received: 181.7 MiB (190500400 bytes)
	Bandwidth per channel: 6.119⇅ Mbps (764.9 kBps)
	Aggregate bandwidth: 152.274↓, 153.698↑ Mbps
	Packet rate estimate: 14586.9↓, 13171.4↑ (2↓, 44↑ TCP MSS/op)
	Test duration: 10.0083 s.
	*/

	/* tcpkali --workers 1 -c 50 -r 5k -m "PING" 192.168.100.120:9000
	Destination: [192.168.100.120]:9000
	Interface eth0 address [192.168.100.1]:0
	Using interface eth0 to connect to [192.168.100.120]:9000
	Ramped up to 50 connections.
	Total data sent:     9.5 MiB (10011720 bytes)
	Total data received: 9.5 MiB (10011720 bytes)
	Bandwidth per channel: 0.320⇅ Mbps (40.0 kBps)
	Aggregate bandwidth: 8.009↓, 8.009↑ Mbps
	Packet rate estimate: 28394.5↓, 28431.6↑ (1↓, 1↑ TCP MSS/op)
	Test duration: 10.0008 s.
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
