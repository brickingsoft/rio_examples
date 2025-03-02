package main

import (
	"flag"
	"fmt"
	echostd "github.com/brickingsoft/rio_examples/benchmark/echo-std"
	"github.com/brickingsoft/rio_examples/benchmark/metric"
	"time"
)

func main() {
	/*
		Port: 9000
		Workers: 5
		Count: 2000
		NBytes: 1024
		ECHO-STD benching complete(1.280013825s): 7940 conn/sec, 7.8M inbounds/sec, 7.8M outbounds/sec, 0 failures
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

	cost, actions, inbounds, outbounds, failures, err = echostd.Bench(workers, count, port, nBytes)
	if err != nil {
		fmt.Println(fmt.Errorf("ECHO-STD benching failed: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("ECHO-STD benching complete(%s): %d conn/sec, %s inbounds/sec, %s outbounds/sec, %d failures",
		cost.String(), actions, metric.FormatBytes(inbounds), metric.FormatBytes(outbounds), failures))
}
