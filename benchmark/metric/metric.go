package metric

import (
	"fmt"
	"sync/atomic"
	"time"
)

func New() *Metric {
	return &Metric{}
}

type Metric struct {
	actions   atomic.Uint64
	inbounds  atomic.Uint64
	outbounds atomic.Uint64
	failures  atomic.Uint64
	beg       time.Time
	end       time.Time
}

func (m *Metric) Begin() {
	m.beg = time.Now()
}

func (m *Metric) End() {
	m.end = time.Now()
}

func (m *Metric) IncACT(n int) {
	m.actions.Add(uint64(n))
}

func (m *Metric) IncIN(n int) {
	m.inbounds.Add(uint64(n))
}

func (m *Metric) IncOUT(n int) {
	m.outbounds.Add(uint64(n))
}

func (m *Metric) Failed(n int) {
	m.failures.Add(uint64(n))
}

func (m *Metric) Failures() uint64 {
	return m.failures.Load()
}

func (m *Metric) CostDuration() time.Duration {
	return m.end.Sub(m.beg)
}

func (m *Metric) PerSecond() (actions uint64, inbounds uint64, outbounds uint64) {
	sec := uint64(m.end.Sub(m.beg).Seconds())
	if sec < 1 {
		sec = 1
	}
	actions = m.actions.Load() / sec
	inbounds = m.inbounds.Load() / sec
	outbounds = m.outbounds.Load() / sec
	return
}

func (m *Metric) String() string {
	return fmt.Sprintf("actions: %d, inbounds: %d, outbounds: %d", m.actions.Load(), m.inbounds.Load(), m.outbounds.Load())
}
