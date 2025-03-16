package main

import (
	"context"
	"crypto/tls"
	"errors"
	"flag"
	"fmt"
	"github.com/aacfactory/afssl"
	"github.com/brickingsoft/rio/security"
	"io"
	"net"
	"strings"
	"time"
)

func main() {
	var port int
	flag.IntVar(&port, "port", 9000, "server port")
	flag.Parse()

	config := afssl.CertificateConfig{}
	// ca
	caPEM, caKeyPEM, caErr := afssl.GenerateCertificate(config, afssl.CA())
	if caErr != nil {
		fmt.Println("caErr:", caErr)
		return
	}
	sc, cc, configErr := afssl.SSC(caPEM, caKeyPEM)
	if configErr != nil {
		fmt.Println("configErr:", configErr)
		return
	}

	ctx, cancel := context.WithCancel(context.Background())
	ln, lnDone, lnErr := listen(ctx, fmt.Sprintf(":%d", port), sc)
	if lnErr != nil {
		fmt.Println("lnErr:", lnErr)
		cancel()
		return
	}
	dialErr := dial(fmt.Sprintf("127.0.0.1:%d", port), cc)
	if dialErr != nil {
		fmt.Println("dialErr:", dialErr)
	}
	cancel()
	ln.Close()
	<-lnDone
}

func listen(ctx context.Context, address string, config *tls.Config) (ln net.Listener, done <-chan struct{}, err error) {
	ln, err = security.Listen("tcp", address, config)
	if err != nil {
		return
	}
	stopCh := make(chan struct{}, 1)
	done = stopCh
	go func(ctx context.Context, ln net.Listener, stopCh chan struct{}) {
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
				if errors.Is(acceptErr, context.DeadlineExceeded) || strings.Contains(acceptErr.Error(), "timeout") {
					break
				}
				fmt.Println("Accept error:", acceptErr)
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
		close(stopCh)
	}(ctx, ln, stopCh)
	return
}

func dial(address string, config *tls.Config) (err error) {
	conn, dialErr := security.Dial("tcp", address, config)
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
