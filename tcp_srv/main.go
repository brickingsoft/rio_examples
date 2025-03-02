package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"log"
	"net"
)

func main() {
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
			b := make([]byte, 1024)
			rn, rErr := conn.Read(b)
			if rErr != nil {
				conn.Close()
				return
			}
			_, wEr := conn.Write(b[:rn])
			if wEr != nil {
				conn.Close()
			}
		}(conn)
	}
}
