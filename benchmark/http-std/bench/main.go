package main

import (
	"flag"
	"fmt"
	httpstd "github.com/brickingsoft/rio_examples/benchmark/http-std"
	"github.com/brickingsoft/rio_examples/benchmark/metric"
	"time"
)

func main() {
	/*
		Port: 9000
		Workers: 5
		Count: 2000
		NBytes: 1024
		HTTP-STD benching complete(758.063536ms): 5124 conn/sec, 5M inbounds/sec, 5M outbounds/sec, 0 failures

		Port: 9000
		Workers: 10
		Count: 5000
		NBytes: 1024
		HTTP-STD benching complete(622.619091ms): 4097 conn/sec, 4M inbounds/sec, 4M outbounds/sec, 0 failures
	*/
	var (
		port    int
		workers int
		count   int
		nBytes  int
	)

	flag.IntVar(&port, "port", 9000, "server port")
	flag.IntVar(&workers, "workers", 5, "workers")
	flag.IntVar(&count, "count", 2000, "count")
	flag.IntVar(&nBytes, "nBytes", 1024, "nBytes")
	flag.Parse()

	var (
		cost      time.Duration
		actions   uint64
		inbounds  uint64
		outbounds uint64
		failures  uint64
		err       error
	)

	// ECHO RIO
	fmt.Println("------ Benchmark ------")
	fmt.Println("Port:", port)
	fmt.Println("Workers:", workers)
	fmt.Println("Count:", count)
	fmt.Println("NBytes:", nBytes)

	cost, actions, inbounds, outbounds, failures, err = httpstd.Bench(workers, count, port, nBytes)
	if err != nil {
		fmt.Println(fmt.Errorf("HTTP-STD benching failed: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("HTTP-STD benching complete(%s): %d conn/sec, %s inbounds/sec, %s outbounds/sec, %d failures",
		cost.String(), actions, metric.FormatBytes(inbounds), metric.FormatBytes(outbounds), failures))
}
