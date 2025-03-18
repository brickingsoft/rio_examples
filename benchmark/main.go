package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio_examples/benchmark/kali"
	"github.com/brickingsoft/rio_examples/benchmark/local"
	"strings"
	"time"
)

func main() {
	var (
		host    string
		port    int
		mode    string
		count   int
		dur     string
		msg     string
		msgSize int
		repeat  int
		out     string
	)
	flag.StringVar(&host, "host", "192.168.100.120", "host for tcpkali mode")
	flag.IntVar(&port, "port", 9000, "server base port")
	flag.StringVar(&mode, "mode", "", "local: bench local case, tcpkali: use tcpkali to bench, server: run tcp server")
	flag.IntVar(&count, "count", 50, "connection count, max is 500")
	flag.StringVar(&dur, "time", "10s", "time duration")
	flag.IntVar(&repeat, "repeat", 0, "repeat per connection")
	flag.StringVar(&out, "out", "", "result output dir")
	flag.IntVar(&msgSize, "msg_size", 0, "message size")
	flag.StringVar(&msg, "msg", "", "message")
	flag.Parse()

	host = strings.TrimSpace(host)
	if port <= 0 || port > 65535 {
		port = 9000
	}
	if out == "" {
		out = "./benchmark/out"
	}
	msg = strings.TrimSpace(msg)
	if msg == "" && msgSize == 0 {
		msg = "PING"
	} else if msg == "" && msgSize > 0 {
		msg = strings.Repeat("A", msgSize)
	}
	if count > 500 {
		count = 50
		fmt.Println("connection count too big, use 50")
	}

	mode = strings.ToLower(strings.TrimSpace(mode))
	switch mode {
	case "tcpkali":
		kali.Bench(host, port, count, repeat, dur, msg, out)
		break
	case "server":
		kali.Serve(port)
		break
	default:
		dur = strings.ToLower(strings.TrimSpace(dur))
		if dur == "" {
			dur = "10s"
		}
		d, dErr := time.ParseDuration(dur)
		if dErr != nil {
			fmt.Println("parse time failed:", dErr)
			fmt.Println("use 10s.")
		}
		if d < 1 {
			d = 10 * time.Second
			fmt.Println("dur is too small, use 10s.")
		}

		local.Bench(port, count, d, msg, out)
		break
	}
}
