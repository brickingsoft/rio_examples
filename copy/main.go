package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/brickingsoft/rio"
	"io"
	"net"
	"sync"
	"time"
)

func main() {
	config := rio.ListenConfig{}
	ln, lnErr := config.Listen(context.Background(), "tcp", "0.0.0.0:8080")
	if lnErr != nil {
		fmt.Println("lnErr:", lnErr)
		return
	}
	go listenUpstream()

	go func() {
		time.Sleep(1 * time.Second)
		conn, connErr := rio.Dial("tcp", "127.0.0.1:8080")
		if connErr != nil {
			fmt.Println("cli > connErr:", connErr)
			return
		}

		for i := 0; i < 5; i++ {
			_, wErr := conn.Write([]byte("hello world"))
			if wErr != nil {
				fmt.Println("cli > write:", wErr)
				return
			}
			time.Sleep(1 * time.Second)
		}

		b := make([]byte, 1024)
		rn, rnErr := conn.Read(b)
		if rnErr != nil {
			fmt.Println("cli > read:", rnErr)
			return
		}
		fmt.Println("cli > read:", string(b[:rn]))
		_ = conn.Close()
	}()

	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Println("srv > accept err:", err)
			return
		}

		go listen(conn)
	}

}

func listen(conn net.Conn) {
	defer conn.Close()

	upstream, err := rio.Dial("tcp", "127.0.0.1:5555") // nc -l 5555
	if err != nil {
		fmt.Println("srv > dial failed:", err)
		return
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func(upstream net.Conn, conn net.Conn, wg *sync.WaitGroup) {
		defer wg.Done()
		written, wtErr := io.Copy(upstream, conn)
		if wtErr != nil {
			fmt.Println("srv > write to failed:", wtErr)
			return
		}
		fmt.Println("srv > write to:", written)
		_ = upstream.Close()
	}(upstream, conn, wg)

	go func(upstream net.Conn, conn net.Conn, wg *sync.WaitGroup) {
		defer wg.Done()
		read, rfErr := io.Copy(conn, upstream)
		if rfErr != nil {
			if errors.Is(rfErr, net.ErrClosed) {
				fmt.Println("srv > read from:", read)
				return
			}
			fmt.Println("srv > read from failed:", rfErr)
			return
		}
		fmt.Println("srv > read from:", read)
	}(upstream, conn, wg)

	wg.Wait()

}

func listenUpstream() {
	ln, lnErr := rio.Listen("tcp", "0.0.0.0:5555")
	if lnErr != nil {
		fmt.Println("upstream > lnErr:", lnErr)
		return
	}

	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return
			}
			fmt.Println("upstream > accept err:", err)
			return
		}

		b := make([]byte, 4096)
		for {
			n, readErr := conn.Read(b)
			if readErr != nil {
				_ = conn.Close()
				if errors.Is(readErr, io.EOF) {
					fmt.Println("upstream > read EOF")
					break
				}
				fmt.Println("upstream > read > error:", readErr)
				return
			}
			fmt.Println("upstream > read >", string(b[:n]))
			_, writeErr := conn.Write(b[:n])
			if writeErr != nil {
				fmt.Println("upstream > write > error:", writeErr)
				return
			}
		}

	}
}
