package main

import (
	"flag"
	"fmt"
	httprio "github.com/brickingsoft/rio_examples/benchmark/http-rio"
	"github.com/brickingsoft/rio_examples/benchmark/metric"
	"time"
)

func main() {
	/*
		Port: 9000
		Workers: 5
		Count: 2000
		NBytes: 1024
		HTTP-RIO benching complete(1.684125038s): 7987 conn/sec, 7.8M inbounds/sec, 7.8M outbounds/sec, 0 failures

		Port: 9000
		Workers: 10
		Count: 5000
		NBytes: 1024
		HTTP-RIO benching complete(1.041162977s): 7372 conn/sec, 7.2M inbounds/sec, 7.2M outbounds/sec, 0 failures
	*/
	var (
		port    int
		workers int
		count   int
		nBytes  int
	)

	flag.IntVar(&port, "port", 9000, "server port")
	flag.IntVar(&workers, "workers", 10, "workers")
	flag.IntVar(&count, "count", 5000, "count")
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

	cost, actions, inbounds, outbounds, failures, err = httprio.Bench(workers, count, port, nBytes)
	if err != nil {
		fmt.Println(fmt.Errorf("HTTP-RIO benching failed: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("HTTP-RIO benching complete(%s): %d conn/sec, %s inbounds/sec, %s outbounds/sec, %d failures",
		cost.String(), actions, metric.FormatBytes(inbounds), metric.FormatBytes(outbounds), failures))
}
