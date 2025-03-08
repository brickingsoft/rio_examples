package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"log"
	"os"
)

func main() {
	/*
		Destination: [192.168.100.120]:9000
		Interface eth0 address [192.168.100.1]:0
		Using interface eth0 to connect to [192.168.100.120]:9000
		Ramped up to 50 connections.
		Total data sent:     261.1 MiB (273747300 bytes)
		Total data received: 259.6 MiB (272187031 bytes)
		Bandwidth per channel: 8.557⇅ Mbps (1069.7 kBps)
		Aggregate bandwidth: 217.585↓, 218.833↑ Mbps
		Packet rate estimate: 25600.4↓, 18948.7↑ (3↓, 30↑ TCP MSS/op)
		Test duration: 10.0075 s.
	*/

	var port int
	flag.IntVar(&port, "port", 9000, "server port")
	flag.Parse()

	os.Setenv("IOURING_USE_CPU_AFFILIATE", "true")
	os.Setenv("IOURING_REG_BUFFERS", "1024,1048576")
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
		go func(conn *rio.TCPConn) {
			buf := conn.AcquireRegisteredBuffer()
			if buf != nil {
				defer conn.ReleaseRegisteredBuffer(buf)
				for {
					_, rErr := conn.ReadFixed(buf)
					if rErr != nil {
						_ = conn.Close()
						return
					}
					_, wEr := conn.WriteFixed(buf)
					if wEr != nil {
						_ = conn.Close()
						return
					}
				}
			} else {
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
			}
		}(conn.(*rio.TCPConn))
	}
}
