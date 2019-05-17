package metric

import (
	"github.com/faycheng/rolling"
	"sync"
	"time"
)

type counter struct {
	tags    []tag
	rolling *rolling.RollingCounter
}

func newCounter() *counter {
	return &counter{}
}

type tag struct {
	key   string
	value string
}

func WithTag(key string, value string) tag {
	return tag{
		key:   key,
		value: value,
	}
}

type Metric struct {
	sync.Mutex
	tags    []tag
	counter rolling.Metric
	gauge   rolling.Metric
}

func NewMetric(tags ...tag) *Metric {
	counterOpts := rolling.RollingCounterOpts{
		Size:           10,
		BucketDuration: time.Millisecond * 100,
	}
	counter := rolling.NewRollingCounter(counterOpts)
	gaugeOpts := rolling.RollingGaugeOpts{
		Size:           10,
		BucketDuration: time.Millisecond * 100,
	}
	gauge := rolling.NewRollingGauge(gaugeOpts)
	return &Metric{
		tags:    tags,
		counter: counter,
		gauge:   gauge,
	}
}

func (m *Metric) Incr() {
	m.Lock()
	defer m.Unlock()
	m.counter.Add(1)
}

func (m *Metric) Gauge(value int64) {
	m.Lock()
	defer m.Unlock()
	m.gauge.Add(value)
}

func (m *Metric) Push() {
	m.Lock()
	defer m.Unlock()
}
