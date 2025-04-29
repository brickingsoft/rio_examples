package srv_gnet

import (
	"bytes"
	"fmt"
	"github.com/panjf2000/gnet/v2"
)

var (
	ShutdownBytes = []byte("<shutdown>")
)

func Serve(port int) {
	go func(port int) {
		echo := &GNetServer{addr: fmt.Sprintf("tcp://:%d", port), multicore: true}
		err := gnet.Run(echo, echo.addr, gnet.WithMulticore(true), gnet.WithLogger(&gnetLogger{}))
		if err != nil {
			panic(err)
		}
	}(port)
	return
}

func ServeAndWait(port int) {
	echo := &GNetServer{addr: fmt.Sprintf("tcp://:%d", port), multicore: true}
	err := gnet.Run(echo, echo.addr, gnet.WithMulticore(true), gnet.WithLogger(&gnetLogger{}))
	if err != nil {
		panic(err)
	}
}

type GNetServer struct {
	gnet.BuiltinEventEngine
	eng       gnet.Engine
	addr      string
	multicore bool
}

func (es *GNetServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	return gnet.None
}

func (es *GNetServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	if bytes.Equal(buf, ShutdownBytes) {
		return gnet.Shutdown
	}
	return gnet.None
}

func (es *GNetServer) OnShutdown(eng gnet.Engine) {
	//fd, fdErr := eng.Dup()
	//if fdErr != nil {
	//	return
	//}
	//fmt.Println("server shutdown", fd)
	//syscall.Close(fd)
}

type gnetLogger struct{}

func (g *gnetLogger) Debugf(format string, args ...any) {
	return
}

func (g *gnetLogger) Infof(format string, args ...any) {
	return
}

func (g *gnetLogger) Warnf(format string, args ...any) {
	return
}

func (g *gnetLogger) Errorf(format string, args ...any) {
	return
}

func (g *gnetLogger) Fatalf(format string, args ...any) {
	return
}
