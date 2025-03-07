package main

import (
	"github.com/brickingsoft/rio_examples/images"
	"strings"
)

func main() {
	values := []float64{27791.8, 22095.3, 14272.9, 15161.3}
	names := []string{"RIO", "GNET", "EVIO", "NET(STD)"}

	out := strings.Replace("out/tcpkali.png", " ", "_", -1)
	images.Plotit(out, "Echo", values, names)
}
