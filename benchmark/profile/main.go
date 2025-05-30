package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/brickingsoft/rio"
	"io"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"sync"
	"time"
)

func main() {
	var port int
	var count int
	var repeat int
	var out string
	flag.IntVar(&port, "port", 9000, "server port")
	flag.IntVar(&count, "count", 50, "count")
	flag.IntVar(&repeat, "repeat", 1000, "repeat")
	flag.StringVar(&out, "out", "", "out directory")
	flag.Parse()

	//rio.Preset(aio.WithWaitCQETimeoutCurve(aio.Curve{{N: 1, Timeout: 50 * time.Microsecond}}))
	//rio.Preset(aio.WithNAPIBusyPollTimeout(50 * time.Microsecond))
	//rio.Preset(aio.WithMultishotDisabled(true))
	//rio.Preset(aio.WithFlags(liburing.IORING_SETUP_SQPOLL | liburing.IORING_SETUP_SQ_AFF | liburing.IORING_SETUP_REGISTERED_FD_ONLY))
	//rio.Preset(aio.WithWaitCQETimeoutCurve(aio.LCurve))
	//rio.Preset(aio.WithFlags(liburing.IORING_SETUP_REGISTERED_FD_ONLY))

	if out == "" {
		out = "./benchmark/out"
	}

	cpuFile, _ := os.Create(filepath.Join(out, "cpu.pprof"))
	defer cpuFile.Close()
	heapFile, _ := os.Create(filepath.Join(out, "heap.pprof"))
	defer heapFile.Close()
	blockFile, _ := os.Create(filepath.Join(out, "block.pprof"))
	defer blockFile.Close()
	goroutineFile, _ := os.Create(filepath.Join(out, "goroutine.pprof"))
	defer goroutineFile.Close()

	// 开始采集 CPU 数据
	pprof.StartCPUProfile(cpuFile)
	defer pprof.StopCPUProfile()

	// 写入当前堆内存快照
	defer pprof.WriteHeapProfile(heapFile)

	// 启用阻塞分析
	runtime.SetBlockProfileRate(1)       // 记录所有阻塞事件
	defer runtime.SetBlockProfileRate(0) // 程序退出前关闭
	defer pprof.Lookup("block").WriteTo(blockFile, 0)

	// 写入 Goroutine 分析数据
	defer pprof.Lookup("goroutine").WriteTo(goroutineFile, 0)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config := rio.ListenConfig{}
	//config := net.ListenConfig{}
	ln, lnErr := config.Listen(ctx, "tcp", fmt.Sprintf(":%d", port))
	if lnErr != nil {
		fmt.Println("lnErr:", lnErr)
		return
	}
	defer ln.Close()
	go listen(ln)

	now := time.Now()
	dial(fmt.Sprintf("127.0.0.1:%d", port), count, repeat)
	fmt.Println("done!!", time.Since(now))
}

func listen(ln net.Listener) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			return
		}
		go func(conn net.Conn) {
			b := make([]byte, 1024)
			for {
				rn, readErr := conn.Read(b)
				if readErr != nil {
					if errors.Is(readErr, io.EOF) {
						break
					}
					fmt.Println("[srv][read]:", readErr)
					break
				}
				_, writeErr := conn.Write(b[:rn])
				if writeErr != nil {
					fmt.Println("[srv][write]:", writeErr)
					break
				}
			}
			_ = conn.Close()
		}(conn)
	}
}

func dial(address string, count int, repeat int) {
	wg := &sync.WaitGroup{}
	wg.Add(count)

	for i := 0; i < count; i++ {
		go func(wg *sync.WaitGroup, address string, repeat int) {
			defer wg.Done()
			dialer := rio.Dialer{
				Timeout:  5 * time.Second,
				Deadline: time.Time{},
			}
			//dialer := net.Dialer{}
			conn, dialErr := dialer.Dial("tcp", address)
			if dialErr != nil {
				fmt.Println("dialErr:", dialErr)
				return
			}
			defer conn.Close()
			b := []byte("hello world")
			for j := 0; j < repeat; j++ {
				_, writeErr := conn.Write(b)
				if writeErr != nil {
					fmt.Println("[cli][write]:", writeErr)
					break
				}
				//fmt.Println("[cli][write]:", wn)
				_, readErr := conn.Read(b)
				if readErr != nil {
					if errors.Is(readErr, io.EOF) {
						break
					}
					fmt.Println("[cli][read]:", readErr)
					break
				}
				//fmt.Println("[cli][read]:", rn)
			}
		}(wg, address, repeat)
	}
	wg.Wait()
	return
}
