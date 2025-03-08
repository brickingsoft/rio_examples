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
		Total data sent:     193.9 MiB (203302768 bytes)
		Total data received: 192.3 MiB (201660064 bytes)
		Bandwidth per channel: 6.476⇅ Mbps (809.5 kBps)
		Aggregate bandwidth: 161.249↓, 162.563↑ Mbps
		Packet rate estimate: 20438.3↓, 14234.6↑ (3↓, 36↑ TCP MSS/op)
		Test duration: 10.0049 s.
	*/

	var port int
	var schema string
	flag.IntVar(&port, "port", 9000, "server port")
	flag.StringVar(&schema, "schema", aio.DefaultFlagsSchema, "iouring schema")
	flag.Parse()

	fmt.Println("schema:", schema)

	switch strings.ToUpper(strings.TrimSpace(schema)) {
	case aio.PerformanceFlagsSchema:
		os.Setenv("IOURING_SETUP_FLAGS_SCHEMA", aio.PerformanceFlagsSchema)
		break
	default:
		break
	}

	os.Setenv("IOURING_USE_CPU_AFFILIATE", "true")
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
