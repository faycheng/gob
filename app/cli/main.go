package main

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/faycheng/gob/bucket"
	"github.com/faycheng/gob/metric"
	"github.com/faycheng/gob/worker"

	"github.com/faycheng/gokit/plugin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	soPlugin   string
	grpcPlugin string
	grpcAddr   string

	duration time.Duration

	constant  bool
	constantC int
	linear    bool
	linearA   float64
	linearB   float64
	power     bool
	powerY    float64
	powerB    float64
	exp2      bool
	exp2B     float64

	influxdb         bool
	influxdbAddr     string
	influxdbUser     string
	influxdbPassword string
	factory          *metric.Factory
	inCounter        *metric.Metric
	passCounter      *metric.Metric
	errCounter       *metric.Metric
	tsGauge          *metric.Metric
)

func initMetric(name string) {
	factory = loadMetricFactory()
	tag := metric.WithTag("gob.name", name)
	inCounter, _ = factory.NewCounter("gob.in", tag)
	passCounter, _ = factory.NewCounter("gob.pass", tag)
	errCounter, _ = factory.NewCounter("gob.err", tag)
	tsGauge, _ = factory.NewGauge("gob.ts", tag)
}

func loadMetricFactory() *metric.Factory {
	if influxdb {
		return metric.NewMetricFactory(metric.WithInfluxDB(influxdbAddr, influxdbUser, influxdbPassword))
	}
	return metric.NewMetricFactory()
}

func loadPlugin() (p plugin.Plugin) {
	if soPlugin != "" {
		p = plugin.NewSoPlugin(soPlugin)
	}
	if grpcPlugin != "" {
		p = plugin.NewGrpcPlugin(grpcPlugin, grpcAddr)
	}
	return p
}

func loadBucket() *bucket.Bucket {
	life := bucket.NewLife(time.Second, time.Now(), time.Now().Add(duration))
	var seq bucket.Sequence
	if constant {
		seq = bucket.NewConstant(constantC, life)
	}
	if linear {
		seq = bucket.NewLinear(linearA, linearB, life)
	}
	if power {
		seq = bucket.NewPower(powerY, powerB, life)
	}
	if exp2 {
		seq = bucket.NewExp2(exp2B, life)
	}
	return bucket.NewBucket(life, seq)
}

func wrapMetric(call plugin.Call) plugin.Call {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		begin := time.Now()
		inCounter.Add(1)
		reply, err = call(ctx, req)
		if err != nil {
			errCounter.Add(1)
			return
		}
		passCounter.Add(1)
		tsGauge.Add(int64(time.Since(begin) / time.Millisecond))
		return reply, err
	}
}

func gob(name string, req interface{}) error {
	initMetric(name)
	plug := loadPlugin()
	bucket := loadBucket()
	pool := worker.NewPool()
	call, err := plug.Lookup(name)
	if err != nil {
		return err
	}
	call = wrapMetric(call)
	for bucket.Get() {
		go func() {
			r := pool.Get()
			// TODO: call timeout
			err := r.Call(call, req)
			if err != nil {
				logrus.Error(err)
				return
			}
		}()
	}
	return nil
}

func main() {
	var gobCmd = &cobra.Command{
		Use: "gob [name] [args...]",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("bad arguments, args:%+v", args)
			}
			var req interface{}
			if len(args) > 1 {
				err := json.Unmarshal([]byte(args[1]), &req)
				if err != nil {
				}
				return errors.Wrapf(err, "bad arguments, args:%+v", args)
			}
			return gob(args[0], req)
		},
	}
	flags := gobCmd.PersistentFlags()
	flags.StringVarP(&soPlugin, "so", "", "", "")
	flags.StringVarP(&grpcPlugin, "grpc", "", "", "")
	flags.StringVarP(&grpcAddr, "grpc.Addr", "", "", "")

	flags.DurationVarP(&duration, "duration", "", time.Second*10, "")

	flags.BoolVarP(&constant, "constant", "", false, "")
	flags.IntVarP(&constantC, "constant.C", "", 100, "")

	flags.BoolVarP(&linear, "linear", "", false, "")
	flags.Float64VarP(&linearA, "linear.A", "", 1, "")
	flags.Float64VarP(&linearB, "linear.B", "", 100, "")

	flags.BoolVarP(&power, "power", "", false, "")
	flags.Float64VarP(&powerY, "power.Y", "", 1, "")
	flags.Float64VarP(&powerB, "power.B", "", 100, "")

	flags.BoolVarP(&exp2, "exp2", "", false, "")
	flags.Float64VarP(&exp2B, "exp2.B", "", 1, "")

	flags.BoolVarP(&influxdb, "influxdb", "", false, "")
	flags.StringVarP(&influxdbAddr, "influxdb.addr", "", "http://127.0.0.1:8086", "")
	flags.StringVarP(&influxdbUser, "influxdb.user", "", "", "")
	flags.StringVarP(&influxdbUser, "influxdb.password", "", "", "")
	if err := gobCmd.Execute(); err != nil {
		panic(err)
	}
}
