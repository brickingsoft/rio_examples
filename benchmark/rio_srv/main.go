package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/iouring/aio"
	"log"
	"net"
)

func main() {

	/* default
	Destination: [192.168.100.120]:9000
	Interface eth0 address [192.168.100.1]:0
	Using interface eth0 to connect to [192.168.100.120]:9000
	Ramped up to 50 connections.
	Total data sent:     368.5 MiB (386415368 bytes)
	Total data received: 366.7 MiB (384519186 bytes)
	Bandwidth per channel: 12.325⇅ Mbps (1540.7 kBps)
	Aggregate bandwidth: 307.375↓, 308.891↑ Mbps
	Packet rate estimate: 35599.0↓, 26694.1↑ (3↓, 26↑ TCP MSS/op)
	Test duration: 10.0078 s.
	*/

	var port int
	var schema string
	flag.IntVar(&port, "port", 9000, "server port")
	flag.StringVar(&schema, "schema", aio.DefaultFlagsSchema, "iouring schema")
	flag.Parse()

	fmt.Println("settings:", port, schema)

	rio.PrepareIOURingSetupOptions(
		aio.WithFlagsSchema(schema),
	)

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
