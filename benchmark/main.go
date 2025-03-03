package main

import (
	"flag"
	"fmt"
	echorio "github.com/brickingsoft/rio_examples/benchmark/echo-rio"
	echostd "github.com/brickingsoft/rio_examples/benchmark/echo-std"
	"github.com/brickingsoft/rio_examples/benchmark/metric"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"strconv"
	"strings"
	"time"
)

func main() {
	var (
		port      int
		workers   int
		count     int
		nBytes    int
		benchHTTP bool
		draw      bool
	)

	flag.IntVar(&port, "port", 9000, "server port")
	flag.IntVar(&workers, "workers", 10, "workers")
	flag.IntVar(&count, "count", 1000, "count")
	flag.IntVar(&nBytes, "nBytes", 1024, "nBytes")
	flag.BoolVar(&benchHTTP, "http", false, "http")
	flag.BoolVar(&draw, "draw", false, "draw")
	flag.Parse()

	var (
		values []float64
		names  []string
		out    string
	)

	var (
		cost      time.Duration
		actions   uint64
		inbounds  uint64
		outbounds uint64
		failures  uint64
		err       error
	)

	fmt.Println("------ Benchmark ------")
	fmt.Println("Port:", port)
	fmt.Println("Workers:", workers)
	fmt.Println("Count:", count)
	fmt.Println("NBytes:", nBytes)

	// ECHO
	names = nil
	values = nil
	port++
	cost, actions, inbounds, outbounds, failures, err = echorio.Bench(workers, count, port, nBytes)
	if err != nil {
		fmt.Println(fmt.Errorf("ECHO-RIO benching failed: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("ECHO-RIO benching complete(%s): %d conn/sec, %s inbounds/sec, %s outbounds/sec, %d failures",
		cost.String(), actions, metric.FormatBytes(inbounds), metric.FormatBytes(outbounds), failures))

	names = append(names, "RIO")
	values = append(values, float64(actions))

	port++
	cost, actions, inbounds, outbounds, failures, err = echostd.Bench(workers, count, port, nBytes)
	if err != nil {
		fmt.Println(fmt.Errorf("ECHO-STD benching failed: %v", err))
		return
	}
	fmt.Println(fmt.Sprintf("ECHO-STD benching complete(%s): %d conn/sec, %s inbounds/sec, %s outbounds/sec, %d failures",
		cost.String(), actions, metric.FormatBytes(inbounds), metric.FormatBytes(outbounds), failures))

	names = append(names, "STD")
	values = append(values, float64(actions))
	if draw {
		out = strings.Replace("out/echo.png", " ", "_", -1)
		plotit(out, "Echo", values, names)
	}

	if benchHTTP {
		// HTTP
		names = nil
		values = nil
		port++
		cost, actions, inbounds, outbounds, failures, err = echorio.Bench(workers, count, port, nBytes)
		if err != nil {
			fmt.Println(fmt.Errorf("HTTP-RIO benching failed: %v", err))
			return
		}
		fmt.Println(fmt.Sprintf("HTTP-RIO benching complete(%s): %d conn/sec, %s inbounds/sec, %s outbounds/sec, %d failures",
			cost.String(), actions, metric.FormatBytes(inbounds), metric.FormatBytes(outbounds), failures))

		names = append(names, "RIO")
		values = append(values, float64(actions))

		port++
		cost, actions, inbounds, outbounds, failures, err = echostd.Bench(workers, count, port, nBytes)
		if err != nil {
			fmt.Println(fmt.Errorf("HTTP-STD benching failed: %v", err))
			return
		}
		fmt.Println(fmt.Sprintf("HTTP-STD benching complete(%s): %d conn/sec, %s inbounds/sec, %s outbounds/sec, %d failures",
			cost.String(), actions, metric.FormatBytes(inbounds), metric.FormatBytes(outbounds), failures))

		names = append(names, "STD")
		values = append(values, float64(actions))
		if draw {
			out = strings.Replace("out/http.png", " ", "_", -1)
			plotit(out, "Http", values, names)
		}
	}

}

func plotit(path, title string, values []float64, names []string) {
	plot.DefaultFont = font.Font{
		Typeface: "Helvetica",
		Variant:  "Serif",
	}
	var groups []plotter.Values
	for _, value := range values {
		groups = append(groups, plotter.Values{value})
	}
	p := plot.New()
	p.Title.Text = title
	p.Y.Tick.Marker = commaTicks{}
	p.Y.Label.Text = "Req/s"
	bw := 25.0
	w := vg.Points(bw)
	var bars []plot.Plotter
	var barsp []*plotter.BarChart
	for i := 0; i < len(values); i++ {
		bar, err := plotter.NewBarChart(groups[i], w)
		if err != nil {
			panic(err)
		}
		bar.LineStyle.Width = vg.Length(0)
		bar.Color = plotutil.Color(i)
		bar.Offset = vg.Length(
			(float64(w) * float64(i)) -
				(float64(w)*float64(len(values)))/2)
		bars = append(bars, bar)
		barsp = append(barsp, bar)
	}
	p.Add(bars...)
	for i, name := range names {
		p.Legend.Add(fmt.Sprintf("%s (%.0f req/s)", name, values[i]), barsp[i])
	}

	p.Legend.Top = true
	p.NominalX("")

	if err := p.Save(7*vg.Inch, 3*vg.Inch, path); err != nil {
		panic(err)
	}
}

type PreciseTicks struct{}

func (PreciseTicks) Ticks(min, max float64) []plot.Tick {
	const suggestedTicks = 3

	if max <= min {
		panic("illegal range")
	}

	tens := math.Pow10(int(math.Floor(math.Log10(max - min))))
	n := (max - min) / tens
	for n < suggestedTicks-1 {
		tens /= 10
		n = (max - min) / tens
	}

	majorMult := int(n / (suggestedTicks - 1))
	switch majorMult {
	case 7:
		majorMult = 6
	case 9:
		majorMult = 8
	}
	majorDelta := float64(majorMult) * tens
	val := math.Floor(min/majorDelta) * majorDelta
	var labels []float64
	for val <= max {
		if val >= min {
			labels = append(labels, val)
		}
		val += majorDelta
	}
	prec := int(math.Ceil(math.Log10(val)) - math.Floor(math.Log10(majorDelta)))
	var ticks []plot.Tick
	for _, v := range labels {
		vRounded := round(v, prec)
		ticks = append(ticks, plot.Tick{Value: vRounded, Label: strconv.FormatFloat(vRounded, 'f', -1, 64)})
	}
	minorDelta := majorDelta / 2
	switch majorMult {
	case 3, 6:
		minorDelta = majorDelta / 3
	case 5:
		minorDelta = majorDelta / 5
	}

	val = math.Floor(min/minorDelta) * minorDelta
	for val <= max {
		found := false
		for _, t := range ticks {
			if t.Value == val {
				found = true
			}
		}
		if val >= min && val <= max && !found {
			ticks = append(ticks, plot.Tick{Value: val})
		}
		val += minorDelta
	}
	return ticks
}

type commaTicks struct{}

func (commaTicks) Ticks(min, max float64) []plot.Tick {
	tks := PreciseTicks{}.Ticks(min, max)
	for i, t := range tks {
		if t.Label == "" {
			continue
		}
		tks[i].Label = addCommas(t.Label)
	}
	return tks
}

func addCommas(s string) string {
	rev := ""
	n := 0
	for i := len(s) - 1; i >= 0; i-- {
		rev += string(s[i])
		n++
		if n%3 == 0 {
			rev += ","
		}
	}
	s = ""
	for i := len(rev) - 1; i >= 0; i-- {
		s += string(rev[i])
	}
	if strings.HasPrefix(s, ",") {
		s = s[1:]
	}
	return s
}

func round(x float64, prec int) float64 {
	if x == 0 {
		return 0
	}
	if prec >= 0 && x == math.Trunc(x) {
		return x
	}
	pow := math.Pow10(prec)
	intermed := x * pow
	if math.IsInf(intermed, 0) {
		return x
	}
	if x < 0 {
		x = math.Ceil(intermed - 0.5)
	} else {
		x = math.Floor(intermed + 0.5)
	}

	if x == 0 {
		return 0
	}
	return x / pow
}
