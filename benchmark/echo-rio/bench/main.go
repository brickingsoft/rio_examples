package main

import (
	"flag"
	"fmt"
	echorio "github.com/brickingsoft/rio_examples/benchmark/echo-rio"
	"github.com/brickingsoft/rio_examples/benchmark/metric"
	"time"
)

func main() {
	/*
		------ Benchmark ------
		Port: 9000
		Workers: 10
		Count: 1000
		NBytes: 1024
		ECHO-RIO benching complete(1.675918093s): 5966 conn/sec, 5.8M inbounds/sec, 5.8M outbounds/sec, 0 failures
	*/
	var (
		port    int
		workers int
		count   int
		nBytes  int
	)

	flag.IntVar(&port, "port", 9000, "server port")
	flag.IntVar(&workers, "workers", 10, "workers")
	flag.IntVar(&count, "count", 1000, "count")
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

	cost, actions, inbounds, outbounds, failures, err = echorio.Bench(workers, count, port, nBytes)
	if err != nil {
		fmt.Println(fmt.Errorf("ECHO-RIO benching failed: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("ECHO-RIO benching complete(%s): %d conn/sec, %s inbounds/sec, %s outbounds/sec, %d failures",
		cost.String(), actions, metric.FormatBytes(inbounds), metric.FormatBytes(outbounds), failures))
}
