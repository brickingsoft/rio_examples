package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/iouring/aio"
	"io"
	"net"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "server port")
	flag.Parse()

	rio.Presets(aio.WithRegisterFixedBuffer(4096, 10))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config := rio.ListenConfig{}
	ln, lnErr := config.Listen(ctx, "tcp", fmt.Sprintf(":%d", port))
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
		fixed, ok := rio.ConvertToFixedReaderWriter(conn)
		if !ok {
			fmt.Println("convert to fixed ReaderWriter")
			return
		}
		buf := fixed.AcquireRegisteredBuffer()
		if buf == nil {
			fmt.Println("no registered buffer")
			return
		}
		rn, readErr := fixed.ReadFixed(buf)
		if readErr != nil {
			fixed.ReleaseRegisteredBuffer(buf)
			_ = conn.Close()
			if errors.Is(readErr, io.EOF) {
				break
			}
			fmt.Println("Srv read error:", readErr)
			break
		}
		b, _ := io.ReadAll(buf)
		fmt.Println("Srv read:", rn, string(b))

		buf.Reset()

		_, _ = buf.Write(b)

		wn, writeErr := fixed.WriteFixed(buf)
		if writeErr != nil {
			fmt.Println("Srv write error:", writeErr)
		} else {
			fmt.Println("Srv write:", wn)
		}
		fixed.ReleaseRegisteredBuffer(buf)
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
