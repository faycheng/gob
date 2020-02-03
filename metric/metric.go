package metric

import (
	"sync"
	"time"

	"github.com/faycheng/rolling"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/sirupsen/logrus"
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

type metricFactoryOpt func(*Factory)

func WithInfluxDB(addr, username, password string) metricFactoryOpt {
	influxClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     addr,
		Username: username,
		Password: password,
	})
	if err != nil {
		panic(err)
	}
	return func(factory *Factory) {
		factory.enableInfluxDB = true
		factory.influxClient = influxClient
	}
}

type rollingMetric interface {
	value() float64
	add(val int64)
}

type counterMetric struct {
	m rolling.Metric
}

func (c *counterMetric) add(val int64) {
	c.m.Add(val)
}

func (c *counterMetric) value() float64 {
	return float64(c.m.Value())
}

type gaugeMetric struct {
	m rolling.Metric
}

func (g *gaugeMetric) add(val int64) {
	g.m.Add(val)
}

func (g *gaugeMetric) value() float64 {
	return g.m.(rolling.Aggregation).Avg()
}

type Factory struct {
	enableInfluxDB bool
	influxClient   client.Client
}

func NewMetricFactory(opts ...metricFactoryOpt) *Factory {
	factory := &Factory{}
	for _, opt := range opts {
		opt(factory)
	}
	return factory
}

func (f *Factory) NewCounter(name string, tags ...tag) (counter *Metric, err error) {
	rolling := rolling.NewRollingCounter(rolling.RollingCounterOpts{
		BucketDuration: time.Millisecond * 10,
		Size:           100,
	})
	m := &counterMetric{m: rolling}
	counter = NewMetric(name, time.Second, m, tags...)
	if f.enableInfluxDB {
		counter.influxClient = f.influxClient
	}
	return
}

func (f *Factory) NewGauge(name string, tags ...tag) (gauge *Metric, err error) {
	rolling := rolling.NewRollingGauge(rolling.RollingGaugeOpts{
		BucketDuration: time.Millisecond * 10,
		Size:           100,
	})
	m := &gaugeMetric{m: rolling}
	gauge = &Metric{
		name:   name,
		tags:   tags,
		metric: m,
	}
	if f.enableInfluxDB {
		gauge.influxClient = f.influxClient
	}
	return
}

type Metric struct {
	sync.Mutex
	influxClient client.Client
	name         string
	tags         []tag
	metric       rollingMetric
	duration     time.Duration
	once         sync.Once
}

func NewMetric(name string, duration time.Duration, metric rollingMetric, tags ...tag) *Metric {
	m := &Metric{
		name:     name,
		tags:     tags,
		metric:   metric,
		duration: duration,
	}
	go m.run()
	return m
}

func (m *Metric) Add(val int64) {
	m.metric.add(val)
}

func (m *Metric) push() {
	points, _ := client.NewBatchPoints(client.BatchPointsConfig{
		Database: "gob",
	})
	point, _ := client.NewPoint(
		m.name,
		map[string]string{},
		map[string]interface{}{"value": m.metric.value()},
		time.Now().Add(-time.Second),
	)
	points.AddPoint(point)
	err := m.influxClient.Write(points)
	if err != nil {
		logrus.Warnf("[metric] failed to flush in-memory points to influxdb database:%s name:%s err:%+v", "gob", m.name, err)
		return
	}
	logrus.Infof("[metric] flush in-memory points to influxdb successfullly  database:%s name:%s", "gob", m.name)
}

func (m *Metric) run() {
	ticker := time.NewTicker(m.duration)
	for range ticker.C {
		if m.influxClient == nil {
			logrus.Infof("[metric] ts:%s name:%s value:%v", time.Now().Add(-time.Second).Format("2006-01-02 15:04:05"), m.name, m.metric.value())
			continue
		}
		go m.push()
	}
}
