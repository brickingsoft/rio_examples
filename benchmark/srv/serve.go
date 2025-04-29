package srv

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
	serve func(port int)
}

var TcpServers = []Server{
	{
		Name:  "RIO",
		serve: srv_rio.Serve,
	},
	{
		Name:  "EVIO",
		serve: srv_evio.Serve,
	},
	{
		Name:  "GNET",
		serve: srv_gnet.Serve,
	},
	{
		Name:  "NET",
		serve: srv_net.Serve,
	},
}

var HttpServers = []Server{
	{
		Name:  "RIO",
		serve: srv_rio.ServeHttp,
	},
	{
		Name:  "NET",
		serve: srv_net.ServeHttp,
	},
}

func ServeTcp(port int) {
	for _, server := range TcpServers {
		port++
		go server.serve(port)
		fmt.Println("[TCP ]["+server.Name+"]", "listen at", port)
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

func ServeHttp(port int) {
	for _, server := range HttpServers {
		port++
		go server.serve(port)
		fmt.Println("[HTTP]["+server.Name+"]", "listen at", port)
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
