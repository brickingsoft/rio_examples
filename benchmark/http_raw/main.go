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
	/* raw

	         /\      Grafana   /‾‾/
	    /\  /  \     |\  __   /  /
	   /  \/    \    | |/ /  /   ‾‾\
	  /          \   |   (  |  (‾)  |
	 / __________ \  |_|\_\  \_____/

	     execution: local
	        script: ./k6.js
	        output: -

	     scenarios: (100.00%) 1 scenario, 100 max VUs, 40s max duration (incl. graceful stop):
	              * default: 100 looping VUs for 10s (gracefulStop: 30s)


	     data_received..................: 94 MB  9.4 MB/s
	     data_sent......................: 107 MB 11 MB/s
	     http_req_blocked...............: avg=895ns   min=342ns   med=469ns    max=8.04ms  p(90)=743ns  p(95)=886ns
	     http_req_connecting............: avg=66ns    min=0s      med=0s       max=1.29ms  p(90)=0s     p(95)=0s
	     http_req_duration..............: avg=1.05ms  min=41.58µs med=892.7µs  max=79.41ms p(90)=1.62ms p(95)=2.25ms
	       { expected_response:true }...: avg=1.05ms  min=41.58µs med=892.7µs  max=79.41ms p(90)=1.62ms p(95)=2.25ms
	     http_req_failed................: 0.00%  0 out of 926061
	     http_req_receiving.............: avg=14.13µs min=2.45µs  med=3.55µs   max=13.21ms p(90)=5.53µs p(95)=9.44µs
	     http_req_sending...............: avg=3.61µs  min=1.06µs  med=1.39µs   max=9.94ms  p(90)=1.97µs p(95)=2.87µs
	     http_req_tls_handshaking.......: avg=0s      min=0s      med=0s       max=0s      p(90)=0s     p(95)=0s
	     http_req_waiting...............: avg=1.04ms  min=35.76µs med=885.8µs  max=79.12ms p(90)=1.6ms  p(95)=2.21ms
	     http_reqs......................: 926061 92595.20398/s
	     iteration_duration.............: avg=1.07ms  min=54.33µs med=905.31µs max=80.66ms p(90)=1.65ms p(95)=2.3ms
	     iterations.....................: 926061 92595.20398/s
	     vus............................: 100    min=100         max=100
	     vus_max........................: 100    min=100         max=100


	running (10.0s), 000/100 VUs, 926061 complete and 0 interrupted iterations
	default ✓ [======================================] 100 VUs  10s
	*/

	/*
			--wait_count=16 --wait_timeout=500us --prepare_timeout=500ns --schema=PERFORMANACE

			execution: local
		        script: ./k6.js
		        output: -

		     scenarios: (100.00%) 1 scenario, 100 max VUs, 40s max duration (incl. graceful stop):
		              * default: 100 looping VUs for 10s (gracefulStop: 30s)


		     data_received..................: 103 MB  10 MB/s
		     data_sent......................: 120 MB  12 MB/s
		     http_req_blocked...............: avg=931ns    min=367ns    med=524ns    max=9.09ms  p(90)=985ns  p(95)=1.46µs
		     http_req_connecting............: avg=46ns     min=0s       med=0s       max=2.12ms  p(90)=0s     p(95)=0s
		     http_req_duration..............: avg=955.16µs min=166.5µs  med=789.98µs max=29.72ms p(90)=1.43ms p(95)=1.9ms
		       { expected_response:true }...: avg=955.16µs min=166.5µs  med=789.98µs max=29.72ms p(90)=1.43ms p(95)=1.9ms
		     http_req_failed................: 0.00%   0 out of 1023631
		     http_req_receiving.............: avg=10.15µs  min=2.67µs   med=4.02µs   max=9.53ms  p(90)=8.59µs p(95)=12.74µs
		     http_req_sending...............: avg=3.38µs   min=1.14µs   med=1.56µs   max=9.94ms  p(90)=3.55µs p(95)=5.13µs
		     http_req_tls_handshaking.......: avg=0s       min=0s       med=0s       max=0s      p(90)=0s     p(95)=0s
		     http_req_waiting...............: avg=941.62µs min=158.5µs  med=782.04µs max=29.7ms  p(90)=1.41ms p(95)=1.88ms
		     http_reqs......................: 1023631 102354.98284/s
		     iteration_duration.............: avg=973.28µs min=178.54µs med=804.27µs max=30.07ms p(90)=1.45ms p(95)=1.94ms
		     iterations.....................: 1023631 102354.98284/s
		     vus............................: 100     min=100          max=100
		     vus_max........................: 100     min=100          max=100

	*/
	var port int
	var schema string
	flag.IntVar(&port, "port", 9000, "server port")
	flag.StringVar(&schema, "schema", aio.DefaultFlagsSchema, "iouring schema")

	flag.Parse()

	fmt.Println("settings:", port, schema)

	rio.PrepareIOURingSetupOptions(
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
