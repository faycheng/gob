package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/faycheng/gob/bucket"
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
	duration   time.Duration
	constant   bool
	constantC  int
	linear     bool
	linearA    float64
	linearB    float64
	power      bool
	powerY     float64
	powerB     float64
	exp2       bool
	exp2B      float64
)

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

func gob(name string, req interface{}) error {
	plug := loadPlugin()
	bucket := loadBucket()
	pool := worker.NewPool()
	call, err := plug.Lookup(name)
	if err != nil {
		return err
	}
	for bucket.Get() {
		go func() {
			r := pool.Get()
			// TODO: call timeout
			err := r.Call(call, req)
			if err != nil {
				logrus.Error(err)
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
	if err := gobCmd.Execute(); err != nil {
		panic(err)
	}
}
