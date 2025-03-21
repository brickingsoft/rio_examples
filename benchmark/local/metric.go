package local

import (
	"bytes"
	"fmt"
	"github.com/brickingsoft/rio_examples/benchmark/commons"
	"sync/atomic"
	"time"
)

func NewMetric(title string, msgSize int) *Metric {
	return &Metric{
		title:   title,
		msgSize: float64(msgSize),
		failed:  atomic.Uint64{},
		rn:      atomic.Uint64{},
		in:      atomic.Uint64{},
		wn:      atomic.Uint64{},
		out:     atomic.Uint64{},
		beg:     time.Time{},
		end:     time.Time{},
	}
}

type Metric struct {
	title   string
	msgSize float64
	failed  atomic.Uint64
	rn      atomic.Uint64
	in      atomic.Uint64
	wn      atomic.Uint64
	out     atomic.Uint64
	beg     time.Time
	end     time.Time
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

func (m *Metric) IncrFailed() {
	m.failed.Add(1)
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

func (m *Metric) Failed() uint64 {
	return m.failed.Load()
}

func (m *Metric) Title() string {
	return m.title
}

func (m *Metric) Rate() float64 {
	d := m.Duration()
	rp := float64(m.TotalReceived()) / d.Seconds() / m.msgSize
	return rp
}

func (m *Metric) String() string {
	buf := bytes.NewBufferString("")
	buf.WriteString("------" + m.title + "------\n")
	buf.WriteString(fmt.Sprintf("Total data sent: %s (%d bytes)\n", commons.FormatBytes(m.TotalSent()), m.TotalSent()))
	buf.WriteString(fmt.Sprintf("Total data received: %s (%d bytes)\n", commons.FormatBytes(m.TotalReceived()), m.TotalReceived()))

	d := m.Duration()
	sp := float64(m.TotalSent()) / d.Seconds() / m.msgSize
	rp := float64(m.TotalReceived()) / d.Seconds() / m.msgSize

	buf.WriteString(fmt.Sprintf("sent/sec: %.2f\n", sp))
	buf.WriteString(fmt.Sprintf("recv/sec: %.2f\n", rp))
	buf.WriteString(fmt.Sprintf("duration: %s\n", d))
	return buf.String()
}
