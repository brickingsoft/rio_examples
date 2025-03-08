package main

import (
	"github.com/brickingsoft/rio_examples/images"
	"strings"
)

func main() {
	values := []float64{20438.3, 18568.5, 17832.6, 14937.1}
	names := []string{"RIO", "EVIO", "GNET", "NET(STD)"}

	out := strings.Replace("out/tcpkali.png", " ", "_", -1)
	images.Plotit(out, "Echo", values, names)
}
