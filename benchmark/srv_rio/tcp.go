package srv_rio

import (
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/liburing/aio"
	"net"
)

func Serve(port int) {
	rio.Preset(
		aio.WithWaitCQETimeoutCurve(aio.LCurve),
	)
	ln, lnErr := rio.Listen("tcp", fmt.Sprintf(":%d", port))
	if lnErr != nil {
		panic(lnErr)
		return
	}
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
						_ = conn.Close()
						return
					}
					_, wEr := conn.Write(packet[:rn])
					if wEr != nil {
						_ = conn.Close()
						return
					}
				}
			}(conn)
		}
	}(ln)
	return
}
