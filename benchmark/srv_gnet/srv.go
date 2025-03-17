package srv_gnet

import (
	"fmt"
	"github.com/panjf2000/gnet/v2"
	"log"
)

func Serve(port int) {
	go func(port int) {
		echo := &gnetServer{addr: fmt.Sprintf("tcp://:%d", port), multicore: true}
		log.Fatal(gnet.Run(echo, echo.addr, gnet.WithMulticore(true), gnet.WithLogger(&gnetLogger{})))
	}(port)
	return
}

type gnetServer struct {
	gnet.BuiltinEventEngine
	eng       gnet.Engine
	addr      string
	multicore bool
}

func (es *gnetServer) OnBoot(eng gnet.Engine) gnet.Action {
	es.eng = eng
	return gnet.None
}

func (es *gnetServer) OnTraffic(c gnet.Conn) gnet.Action {
	buf, _ := c.Next(-1)
	c.Write(buf)
	return gnet.None
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
