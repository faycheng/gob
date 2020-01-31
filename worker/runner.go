package worker

import (
	"time"

	"context"
	"fmt"

	"github.com/faycheng/gokit/plugin"
	"github.com/pkg/errors"
)

type options struct {
	Timeout time.Duration
}

func defaultOptions() *options {
	return &options{Timeout: time.Second * 10}
}

type CallOption func(opts *options)

func Timeout(d time.Duration) CallOption {
	return func(opts *options) {
		opts.Timeout = d
	}
}

type Runner interface {
	Call(call plugin.Call, req interface{}, options ...CallOption) error
}

type runner struct {
}

func (r *runner) Call(call plugin.Call, req interface{}, options ...CallOption) error {
	callDone := make(chan error)
	callOptions := defaultOptions()
	for _, opt := range options {
		opt(callOptions)
	}
	ctx, cancel := context.WithTimeout(context.TODO(), callOptions.Timeout)
	go func() {
		// must close callDone after send err to channel
		defer close(callDone)
		var err error
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("failed to call handler, %+v", r)
				err = errors.WithStack(err)
			}
			callDone <- err
		}()
		_, err = call(ctx, req)
	}()
	var err error
	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			err = ctx.Err()
		case err = <-callDone:
		}
	}()
	return err
}
