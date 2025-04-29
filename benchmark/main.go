package main

import (
	"flag"
	"github.com/brickingsoft/rio_examples/benchmark/cli/k6"
	"github.com/brickingsoft/rio_examples/benchmark/cli/kali"
	"github.com/brickingsoft/rio_examples/benchmark/srv"
	"strings"
)

func main() {
	var (
		host string
		port int
		kind string
		mode string
		out  string
	)
	flag.StringVar(&host, "srv-host", "192.168.100.120", "server host")
	flag.IntVar(&port, "port", 9000, "server base port")
	flag.StringVar(&kind, "kind", "http", "tcp or http")
	flag.StringVar(&mode, "mode", "client", "client: use tcpkali and wrk to bench, server: run tcp and http server")
	flag.StringVar(&out, "out", "", "result output dir")
	flag.Parse()

	host = strings.TrimSpace(host)
	if port <= 0 || port > 65535 {
		port = 9000
	}
	if out == "" {
		out = "./benchmark/out"
	}

	kind = strings.ToLower(strings.TrimSpace(kind))
	mode = strings.ToLower(strings.TrimSpace(mode))
	switch kind {
	case "http":
		switch mode {
		case "client":
			// k6
			k6.Bench(host, port, 100, "10s", out)
			break
		default:
			srv.ServeHttp(port)
			break
		}
		break
	default:
		switch mode {
		case "client":
			// C50T10
			kali.Bench(host, port, 50, 0, "10s", out)
			// C50R5K
			kali.Bench(host, port, 50, 5000, "", out)
			break
		default:
			srv.ServeTcp(port)
			break
		}
	}

}
