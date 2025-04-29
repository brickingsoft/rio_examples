package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"io"
	"net"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "server port")
	flag.Parse()

	ln, lnErr := rio.Listen("tcp", fmt.Sprintf(":%d", port))
	if lnErr != nil {
		fmt.Println("lnErr:", lnErr)
		return
	}
	defer ln.Close()
	go listen(ln)

	dialErr := dial(fmt.Sprintf("127.0.0.1:%d", port))
	if dialErr != nil {
		fmt.Println("dialErr:", dialErr)
	}

}

func listen(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Println("accept err:", err)
			return
		}
		b := make([]byte, 1024)
		rn, readErr := conn.Read(b)
		if readErr != nil {
			_ = conn.Close()
			if errors.Is(readErr, io.EOF) {
				break
			}
			fmt.Println("Srv read error:", readErr)
			break
		}
		fmt.Println("Srv read:", rn, string(b[:rn]))

		wn, writeErr := conn.Write(b[:rn])
		if writeErr != nil {
			fmt.Println("Srv write error:", writeErr)
		} else {
			fmt.Println("Srv write:", wn)
		}
		_ = conn.Close()
	}
	return
}

func dial(address string) (err error) {
	conn, dialErr := rio.DialTimeout("tcp", address, 5*time.Second)
	if dialErr != nil {
		err = dialErr
		return
	}
	defer conn.Close()

	wn, writeErr := conn.Write([]byte("hello world"))
	if writeErr != nil {
		fmt.Println("Cli write error:", writeErr)
		return
	}
	fmt.Println("Cli write:", wn)

	b := make([]byte, 1024)
	rn, readErr := conn.Read(b)
	if readErr != nil {
		fmt.Println("Cli read error:", readErr)
		return
	}
	fmt.Println("Cli read:", rn, string(b[:rn]))
	return
}
