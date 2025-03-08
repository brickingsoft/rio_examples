package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/iouring/aio"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	/*
		Destination: [192.168.100.120]:9000
		Interface eth0 address [192.168.100.1]:0
		Using interface eth0 to connect to [192.168.100.120]:9000
		Ramped up to 50 connections.
		Total data sent:     265.5 MiB (278387192 bytes)
		Total data received: 264.1 MiB (276937819 bytes)
		Bandwidth per channel: 8.879⇅ Mbps (1109.8 kBps)
		Aggregate bandwidth: 221.386↓, 222.544↑ Mbps
		Packet rate estimate: 25033.5↓, 19456.6↑ (3↓, 31↑ TCP MSS/op)
		Test duration: 10.0074 s.
	*/

	var port int
	var schema string
	flag.IntVar(&port, "port", 9000, "server port")
	flag.StringVar(&schema, "schema", aio.DefaultFlagsSchema, "iouring schema")
	flag.Parse()

	switch strings.ToUpper(strings.TrimSpace(schema)) {
	case aio.PerformanceFlagsSchema:
		os.Setenv("IOURING_SETUP_FLAGS_SCHEMA", aio.PerformanceFlagsSchema)
		//os.Setenv("IOURING_ENTRIES", "4096")
		break
	default:
		break
	}

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
