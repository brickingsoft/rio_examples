package srv_net

import (
	"fmt"
	"net"
	"time"
)

func Serve(port int) {
	ln, lnErr := net.Listen("tcp", fmt.Sprintf(":%d", port))
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
