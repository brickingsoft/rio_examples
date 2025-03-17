package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/iouring/aio"
	"log"
	"net"
)

func main() {

	/* tcpkali --workers 1 -c 50 -T 10s -m "PING" 192.168.100.120:9000
	Destination: [192.168.100.120]:9000
	Interface eth0 address [192.168.100.1]:0
	Using interface eth0 to connect to [192.168.100.120]:9000
	Ramped up to 50 connections.
	Total data sent:     245.0 MiB (256897392 bytes)
	Total data received: 244.2 MiB (256063278 bytes)
	Bandwidth per channel: 8.202⇅ Mbps (1025.3 kBps)
	Aggregate bandwidth: 204.723↓, 205.390↑ Mbps
	Packet rate estimate: 24043.6↓, 17801.4↑ (3↓, 25↑ TCP MSS/op)
	Test duration: 10.0062 s.
	*/

	/* tcpkali --workers 1 -c 50 -r 5k -m "PING" 192.168.100.120:9000
	Destination: [192.168.100.120]:9000
	Interface eth0 address [192.168.100.1]:0
	Using interface eth0 to connect to [192.168.100.120]:9000
	Ramped up to 50 connections.
	Total data sent:     9.6 MiB (10019512 bytes)
	Total data received: 9.6 MiB (10019132 bytes)
	Bandwidth per channel: 0.320⇅ Mbps (40.0 kBps)
	Aggregate bandwidth: 8.010↓, 8.010↑ Mbps
	Packet rate estimate: 44138.9↓, 44153.1↑ (1↓, 1↑ TCP MSS/op)
	Test duration: 10.0069 s.
	*/

	var port int
	var schema string
	var fixedFiles int
	var autoInstall bool
	var multiAccept bool
	var reusePort bool
	flag.IntVar(&port, "port", 9000, "server port")
	flag.IntVar(&fixedFiles, "files", 1024, "fixed files")
	flag.BoolVar(&autoInstall, "auto", false, "auto install fixed fd")
	flag.BoolVar(&multiAccept, "ma", false, "multi-accept")
	flag.BoolVar(&reusePort, "reuse", false, "reuse port")
	flag.StringVar(&schema, "schema", aio.DefaultFlagsSchema, "iouring schema")
	flag.Parse()

	fmt.Println("settings:", port, schema)

	rio.Presets(
		aio.WithFlagsSchema(schema),
		aio.WithRegisterFixedFiles(uint32(fixedFiles)),
	)

	config := rio.ListenConfig{
		FastOpen:           true,
		QuickAck:           true,
		ReusePort:          reusePort,
		SendZC:             false,
		MultishotAccept:    multiAccept,
		AutoFixedFdInstall: autoInstall,
	}
	ctx := context.Background()
	ln, lnErr := config.Listen(ctx, "tcp", fmt.Sprintf(":%d", port))
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
