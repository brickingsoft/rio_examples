package main

import (
	"github.com/brickingsoft/rio_examples/images"
	"strings"
)

func main() {
	values := []float64{25033.5, 18635.1, 19344.9, 14937.1}
	names := []string{"RIO", "GNET", "EVIO", "NET(STD)"}

	out := strings.Replace("out/tcpkali.png", " ", "_", -1)
	images.Plotit(out, "Echo", values, names)
}
