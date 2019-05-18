package metric

import (
	"sync"
	"time"

	"github.com/faycheng/rolling"
	"github.com/influxdata/influxdb/client/v2"
)

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

type metricFactoryOpt func(*metricFactory)

func WithInfluxDB(addr, username, password string) metricFactoryOpt {
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	return func(factory *metricFactory) {
		factory.enableInfluxDB = true
		factory.influxClient = influxClient
	}

}

type metricFactory struct {
	enableInfluxDB bool
	influxClient   client.Client
}

func NewMetricFacotry(opts ...metricFactoryOpt) *metricFactory {
	factory := &metricFactory{}
	for _, opt := range opts {
		opt(factory)
	}
	return factory
}

func (f *metricFactory) NewCounter(tags ...tag) (counter *metric, err error) {
	timingCounter := rolling.NewTimingCounter(rolling.TimingCounterOpts{
		Size:           3600,
		BucketDuration: time.Second,
	})
	counter = &metric{
		tags:   tags,
		metric: timingCounter,
	}

	return
}

func (f *metricFactory) NewGauge(tags ...tag) (gauge *metric, err error) {
	timingCounter := rolling.NewTimingGauge(rolling.TimingGaugeOpts{
		Size:           3600,
		BucketDuration: time.Second,
	})
	gauge = &metric{
		tags:   tags,
		metric: timingCounter,
	}
	if f.enableInfluxDB {
		gauge.influxClient = f.influxClient
	}
	return

}

type metric struct {
	sync.Mutex
	influxClient client.Client
	tags         []tag
	metric       rolling.Metric
}

func (m *metric) Add(val int64) {
	m.metric.Add(val)
}
