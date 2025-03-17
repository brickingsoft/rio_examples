package kali

import (
	"fmt"
	"github.com/brickingsoft/rio_examples/benchmark/srv_evio"
	"github.com/brickingsoft/rio_examples/benchmark/srv_gnet"
	"github.com/brickingsoft/rio_examples/benchmark/srv_net"
	"github.com/brickingsoft/rio_examples/benchmark/srv_rio"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Name  string
	Serve func(port int)
}

var servers = []Server{
	{
		Name:  "RIO",
		Serve: srv_rio.Serve,
	},
	{
		Name:  "EVIO",
		Serve: srv_evio.Serve,
	},
	{
		Name:  "GNET",
		Serve: srv_gnet.Serve,
	},
	{
		Name:  "NET",
		Serve: srv_net.Serve,
	},
}

func Serve(port int) {
	for _, server := range servers {
		port++
		go server.Serve(port)
		fmt.Println("["+server.Name+"]", "listen at", port)
		time.Sleep(time.Millisecond * 500)
	}
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh,
		syscall.SIGINT,
		syscall.SIGKILL,
		syscall.SIGQUIT,
		syscall.SIGABRT,
		syscall.SIGTERM,
	)
	<-signalCh
}
