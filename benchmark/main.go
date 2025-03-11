package main

import (
	"github.com/brickingsoft/rio_examples/images"
	"strings"
)

func main() {
	values := []float64{35599.0, 18568.5, 17832.6, 14937.1}
	names := []string{"RIO(DEFAULT)", "EVIO", "GNET", "NET"}

	out := strings.Replace("out/bench_echo.png", " ", "_", -1)
	images.Plotit(out, "Echo", values, names)
}
