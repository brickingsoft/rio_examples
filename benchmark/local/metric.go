package local

import (
	"bytes"
	"fmt"
	"github.com/brickingsoft/rio_examples/benchmark/commons"
	"sync/atomic"
	"time"
)

func NewMetric(title string) *Metric {
	return &Metric{
		title: title,
	}
}

type Metric struct {
	title string
	rn    atomic.Uint64
	in    atomic.Uint64
	wn    atomic.Uint64
	out   atomic.Uint64
	beg   time.Time
	end   time.Time
}

func (m *Metric) Begin() {
	m.beg = time.Now()
}

func (m *Metric) End() {
	m.end = time.Now()
}

func (m *Metric) IncrIn(n int) {
	m.in.Add(uint64(n))
	m.rn.Add(1)
}

func (m *Metric) IncrOut(n int) {
	m.out.Add(uint64(n))
	m.wn.Add(1)
}

func (m *Metric) TotalSent() uint64 {
	return m.out.Load()
}

func (m *Metric) TotalReceived() uint64 {
	return m.in.Load()
}

func (m *Metric) Duration() time.Duration {
	return m.end.Sub(m.beg)
}

func (m *Metric) Title() string {
	return m.title
}

func (m *Metric) Rate() float64 {
	d := m.Duration()
	rp := float64(m.TotalReceived()) / d.Seconds()
	return rp
}

func (m *Metric) String() string {
	buf := bytes.NewBufferString("")
	buf.WriteString("------" + m.title + "------\n")
	buf.WriteString(fmt.Sprintf("Total data sent: %s (%d bytes)\n", commons.FormatBytes(m.TotalSent()), m.TotalSent()))
	buf.WriteString(fmt.Sprintf("Total data received: %s (%d bytes)\n", commons.FormatBytes(m.TotalReceived()), m.TotalReceived()))

	d := m.Duration()
	sp := float64(m.TotalSent()) / d.Seconds()
	rp := float64(m.TotalReceived()) / d.Seconds()

	buf.WriteString(fmt.Sprintf("sent/sec: %.2f\n", sp))
	buf.WriteString(fmt.Sprintf("recv/sec: %.2f\n", rp))
	return buf.String()
}
