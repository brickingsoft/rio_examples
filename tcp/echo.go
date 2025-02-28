package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"io"
	"net"
	"strings"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "server port")
	flag.Parse()

	setting()
	ctx, cancel := context.WithCancel(context.Background())
	lnDone, lnErr := listen(ctx, fmt.Sprintf(":%d", port))
	if lnErr != nil {
		fmt.Println("lnErr:", lnErr)
		cancel()
		return
	}
	dialErr := dial(fmt.Sprintf("127.0.0.1:%d", port))
	if dialErr != nil {
		fmt.Println("dialErr:", dialErr)
	}
	cancel()
	<-lnDone
}

func listen(ctx context.Context, address string) (done <-chan struct{}, err error) {
	ln, lnErr := rio.Listen("tcp", address)
	if lnErr != nil {
		err = lnErr
		return
	}
	stopCh := make(chan struct{}, 1)
	done = stopCh
	go func(ctx context.Context, ln net.Listener, stopCh chan struct{}) {
		tcpLn := ln.(*rio.TCPListener)
		stopped := false
		for {
			select {
			case <-ctx.Done():
				stopped = true
				break
			default:
				_ = tcpLn.SetDeadline(time.Now().Add(1 * time.Second))
				conn, acceptErr := ln.Accept()
				if acceptErr != nil {
					if errors.Is(acceptErr, context.DeadlineExceeded) || strings.Contains(acceptErr.Error(), "timeout") {
						break
					}
					fmt.Println("Accept error:", acceptErr)
					stopped = true
					break
				}
				_ = conn.SetReadDeadline(time.Now().Add(1 * time.Second))

				b := make([]byte, 1024)
				rn, readErr := conn.Read(b)
				if readErr != nil {
					_ = conn.Close()
					if errors.Is(err, io.EOF) {
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
				break
			}
			if stopped {
				break
			}
		}
		_ = ln.Close()
		close(stopCh)
	}(ctx, ln, stopCh)
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
