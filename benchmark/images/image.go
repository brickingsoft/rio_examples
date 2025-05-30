package images

import (
	"fmt"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"math"
	"slices"
	"strconv"
	"strings"
)

type Request struct {
	Path  string
	Title string
	Label string // "req/s"
	Items []Item
}

type Item struct {
	Name  string
	Value float64
}

func Draw(req Request) error {
	plot.DefaultFont = font.Font{
		Typeface: "Helvetica",
		Variant:  "Serif",
	}

	slices.SortFunc[[]Item](req.Items, func(a, b Item) int {
		n := a.Value - b.Value
		if n < 0 {
			return -1
		}
		if n > 0 {
			return 1
		}
		return 0
	})

	slices.Reverse[[]Item](req.Items)

	var groups []plotter.Values
	for _, item := range req.Items {
		groups = append(groups, plotter.Values{item.Value})
	}
	p := plot.New()
	p.Title.Text = req.Title
	p.Y.Tick.Marker = commaTicks{}
	p.Y.Label.Text = req.Label
	bw := 25.0
	w := vg.Points(bw)
	var bars []plot.Plotter
	var barsp []*plotter.BarChart
	for i := 0; i < len(req.Items); i++ {
		bar, err := plotter.NewBarChart(groups[i], w)
		if err != nil {
			return err
		}
		bar.LineStyle.Width = vg.Length(0)
		bar.Color = plotutil.Color(i)
		bar.Offset = vg.Length(
			(float64(w) * float64(i)) -
				(float64(w)*float64(len(req.Items)))/2)
		bars = append(bars, bar)
		barsp = append(barsp, bar)
	}
	p.Add(bars...)

	// get max
	maxValue := req.Items[0].Value
	for i := 1; i < len(req.Items); i++ {
		target := req.Items[i].Value
		if target > maxValue {
			maxValue = target
		}
	}
	for i, item := range req.Items {
		x := (item.Value / maxValue) * 100
		p.Legend.Add(fmt.Sprintf("%s (%.0f %%) (%.0f req/s)", item.Name, x, item.Value), barsp[i])
	}

	p.Legend.Top = true
	p.NominalX("")

	if err := p.Save(7*vg.Inch, 3*vg.Inch, req.Path); err != nil {
		return err
	}
	return nil
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
