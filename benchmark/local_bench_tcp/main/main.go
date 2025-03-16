package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio_examples/benchmark/local_bench_tcp"
	"io"
	"net"
	"time"
)

func main() {
	rm, rmErr := local_bench_tcp.Bench2("ECHO", 9001, 50, 10*time.Second, serveRIO)
	if rmErr != nil {
		fmt.Println("err:", rmErr)
		return
	}
	fmt.Println(rm)
}

func main2() {
	var address string
	var count int
	var dur string
	flag.StringVar(&address, "port", "192.168.100.120:9000", "server address")
	flag.IntVar(&count, "count", 50, "connection count")
	flag.StringVar(&dur, "time", "10s", "time duration")
	flag.Parse()

	d, dErr := time.ParseDuration(dur)
	if dErr != nil {
		d = time.Second * 10
	}

	rm, rmErr := local_bench_tcp.Bench("ECHO", address, count, d)
	if rmErr != nil {
		fmt.Println("err:", rmErr)
		return
	}
	fmt.Println(rm)

	/* rio
	Total data sent: 5.3M (5605276 bytes)
	Total data received: 5.3M (5605276 bytes)
	sent/sec: 560314.05
	recv/sec: 560314.05
	*/

	/* EVIO
	Total data sent: 6.3M (6577604 bytes)
	Total data received: 6.3M (6577604 bytes)
	sent/sec: 656747.28
	recv/sec: 656747.28
	*/

	/* GNET
	Total data sent: 6.7M (7015848 bytes)
	Total data received: 6.7M (7015848 bytes)
	sent/sec: 701354.37
	recv/sec: 701354.37
	*/

	//nm, nmErr := local_bench_tcp.Bench("NET", address, count, d)
	//if nmErr != nil {
	//	fmt.Println("NET err:", nmErr)
	//	return
	//}
	//fmt.Println(nm)
}

func serveRIO(port int) (closer io.Closer, err error) {
	ln, lnErr := rio.Listen("tcp", fmt.Sprintf(":%d", port))
	if lnErr != nil {
		err = lnErr
		return
	}
	closer = ln
	go func(ln net.Listener) {
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
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
	}(ln)
	return
}

func serveNet(port int) (closer io.Closer, err error) {
	ln, lnErr := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if lnErr != nil {
		err = lnErr
		return
	}
	closer = ln
	go func(ln net.Listener) {
		for {
			conn, acceptErr := ln.Accept()
			if acceptErr != nil {
				return
			}
			go func(conn net.Conn) {
				var packet [0xFFF]byte
				for {
					conn.SetDeadline(time.Now().Add(15 * time.Second))
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
	}(ln)
	return
}
