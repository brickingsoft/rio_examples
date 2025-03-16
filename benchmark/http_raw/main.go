package main

import (
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"github.com/brickingsoft/rio/pkg/iouring/aio"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	/*  wrk -t 10 -c 1000 -d 10s http://192.168.100.120:9000/
	Running 10s test @ http://192.168.100.120:9000/
	  10 threads and 1000 connections
	  Thread Stats   Avg      Stdev     Max   +/- Stdev
	    Latency    15.50ms   50.87ms 877.13ms   94.65%
	    Req/Sec    19.61k     7.53k   48.56k    68.57%
	  1941877 requests in 10.08s, 187.04MB read
	Requests/sec: 192568.74
	Transfer/sec:     18.55MB
	*/

	/* net
	Running 10s test @ http://192.168.100.120:9000/
	  10 threads and 1000 connections
	  Thread Stats   Avg      Stdev     Max   +/- Stdev
	    Latency    24.91ms   84.51ms 997.16ms   94.28%
	    Req/Sec    18.51k     5.65k   58.13k    66.36%
	  1836687 requests in 10.09s, 222.45MB read
	Requests/sec: 181953.24
	Transfer/sec:     22.04MB
	*/

	/* evio
	Running 10s test @ http://192.168.100.120:9000/
	  10 threads and 1000 connections
	  Thread Stats   Avg      Stdev     Max   +/- Stdev
	    Latency    14.03ms   41.06ms 469.88ms   95.05%
	    Req/Sec    18.75k     6.37k   54.72k    59.35%
	  1860250 requests in 10.10s, 184.50MB read
	Requests/sec: 184222.91
	Transfer/sec:     18.27MB
	*/

	var port int
	var schema string
	var fixedFiles int
	var autoInstall bool
	var multiAccept bool
	var reusePort bool
	flag.IntVar(&port, "port", 9000, "server port")
	flag.IntVar(&fixedFiles, "files", 9000, "fixed files")
	flag.BoolVar(&autoInstall, "auto", false, "auto install fixed fd")
	flag.BoolVar(&multiAccept, "ma", false, "multi-accept")
	flag.BoolVar(&reusePort, "reuse", false, "reuse port")
	flag.StringVar(&schema, "schema", aio.DefaultFlagsSchema, "iouring schema")
	flag.Parse()

	flag.Parse()

	fmt.Println("settings:", port, schema)

	rio.Presets(
		aio.WithFlagsSchema(schema),
	)

	ln, lnErr := rio.Listen("tcp", fmt.Sprintf(":%d", port))
	if lnErr != nil {
		log.Fatal("lnErr:", lnErr)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go func(conn net.Conn) {
			var packet [0xFFF]byte
			is := &InputStream{}
			res := []byte("hello world")
			for {
				rn, rErr := conn.Read(packet[:])
				if rErr != nil {
					conn.Close()
					return
				}
				data := is.Begin(packet[:rn])
				// process the pipeline
				var req request
				out := make([]byte, 0, 1)
				closed := false
				for {
					leftover, reqErr := parsereq(data, &req)
					if reqErr != nil {
						// bad thing happened
						out = appendresp(out, []byte("500 Error"), []byte(reqErr.Error()+"\n"))
						closed = true
						break
					} else if len(leftover) == len(data) {
						// request not ready, yet
						break
					}
					// handle the request
					req.remoteAddr = conn.RemoteAddr().String()
					out = appendhandle(out, res)
					data = leftover
				}

				is.End(data)
				if len(out) > 0 {
					_, wEr := conn.Write(out)
					if wEr != nil {
						closed = true
					}
				}
				if closed {
					_ = conn.Close()
				}
			}
		}(conn)
	}
}

type InputStream struct{ b []byte }

// Begin accepts a new packet and returns a working sequence of
// unprocessed bytes.
func (is *InputStream) Begin(packet []byte) (data []byte) {
	data = packet
	if len(is.b) > 0 {
		is.b = append(is.b, data...)
		data = is.b
	}
	return data
}

// End shifts the stream to match the unprocessed data.
func (is *InputStream) End(data []byte) {
	if len(data) > 0 {
		if len(data) != len(is.b) {
			is.b = append(is.b[:0], data...)
		}
	} else if len(is.b) > 0 {
		is.b = is.b[:0]
	}
}

type request struct {
	proto, method string
	path, query   string
	head, body    string
	remoteAddr    string
}

// appendhandle handles the incoming request and appends the response to
// the provided bytes, which is then returned to the caller.
func appendhandle(b []byte, body []byte) []byte {
	return appendresp(b, []byte("200 OK"), body)
}

// appendresp will append a valid http response to the provide bytes.
// The status param should be the code plus text such as "200 OK".
// The head parameter should be a series of lines ending with "\r\n" or empty.
func appendresp(b []byte, status, body []byte) []byte {
	b = append(b, "HTTP/1.1"...)
	b = append(b, ' ')
	b = append(b, status...)
	b = append(b, '\r', '\n')
	b = append(b, "Server: evio\r\n"...)
	b = append(b, "Date: "...)
	b = time.Now().AppendFormat(b, "Mon, 02 Jan 2006 15:04:05 GMT")
	b = append(b, '\r', '\n')
	if len(body) > 0 {
		b = append(b, "Content-Length: "...)
		b = strconv.AppendInt(b, int64(len(body)), 10)
		b = append(b, '\r', '\n')
	}
	b = append(b, '\r', '\n')
	if len(body) > 0 {
		b = append(b, body...)
	}
	return b
}

// parsereq is a very simple http request parser. This operation
// waits for the entire payload to be buffered before returning a
// valid request.
func parsereq(data []byte, req *request) (leftover []byte, err error) {
	sdata := string(data)
	var i, s int
	var top string
	var clen int
	var q = -1
	// method, path, proto line
	for ; i < len(sdata); i++ {
		if sdata[i] == ' ' {
			req.method = sdata[s:i]
			for i, s = i+1, i+1; i < len(sdata); i++ {
				if sdata[i] == '?' && q == -1 {
					q = i - s
				} else if sdata[i] == ' ' {
					if q != -1 {
						req.path = sdata[s:q]
						req.query = req.path[q+1 : i]
					} else {
						req.path = sdata[s:i]
					}
					for i, s = i+1, i+1; i < len(sdata); i++ {
						if sdata[i] == '\n' && sdata[i-1] == '\r' {
							req.proto = sdata[s:i]
							i, s = i+1, i+1
							break
						}
					}
					break
				}
			}
			break
		}
	}
	if req.proto == "" {
		return data, fmt.Errorf("malformed request")
	}
	top = sdata[:s]
	for ; i < len(sdata); i++ {
		if i > 1 && sdata[i] == '\n' && sdata[i-1] == '\r' {
			line := sdata[s : i-1]
			s = i + 1
			if line == "" {
				req.head = sdata[len(top)+2 : i+1]
				i++
				if clen > 0 {
					if len(sdata[i:]) < clen {
						break
					}
					req.body = sdata[i : i+clen]
					i += clen
				}
				return data[i:], nil
			}
			if strings.HasPrefix(line, "Content-Length:") {
				n, err := strconv.ParseInt(strings.TrimSpace(line[len("Content-Length:"):]), 10, 64)
				if err == nil {
					clen = int(n)
				}
			}
		}
	}
	// not enough data
	return data, nil
}
