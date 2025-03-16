package main

import (
	"github.com/brickingsoft/rio_examples/images"
	"strings"
)

func main() {
	var (
		values []float64
		names  []string
		out    string
	)

	values = []float64{24043.6, 19010.4, 18598.8, 14586.9}
	names = []string{"RIO", "EVIO", "GNET", "NET"}

	out = strings.Replace("benchmark/out/bench_c50t10s.png", " ", "_", -1)
	images.Plotit(out, "Echo(C50 T10s)", values, names)

	values = []float64{44138.9, 29327.7, 28936.6, 28394.5}
	names = []string{"RIO", "EVIO", "GNET", "NET"}

	out = strings.Replace("benchmark/out/bench_c50r5k.png", " ", "_", -1)
	images.Plotit(out, "Echo(C50 R5K)", values, names)
}
