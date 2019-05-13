package worker

import (
	"fmt"
	"time"

	"github.com/faycheng/gob/bucket"
	"github.com/faycheng/gob/task"
)

const (
	WorkerModeQps         = "qps"
	WorkerModeConcurrency = "concurrency"

	BucketModeConstant = "constant"
	BucketModeUp       = "up"
	BucketModelDown    = "down"
	BucketModeRange    = "range"
)

type workerModeNotSupportErr struct{}

func (e *workerModeNotSupportErr) Error() string {
	return fmt.Sprintf(
		"Worker mode is not supported, expect modes: %+v",
		[]string{WorkerModeQps, WorkerModeConcurrency},
	)
}

type workerFactory struct{}

type workerOpts struct {
	mode        string
	bucket      string
	rate        int
	concurrency int
	duration    time.Duration
}

type WorkerOpt interface {
	Apply(*workerOpts)
}

func defaultOpts() *workerOpts {
	return &workerOpts{
		mode:     WorkerModeQps,
		rate:     100,
		duration: time.Second * 10,
	}
}

type funcWorkerOpt struct {
	f func(*workerOpts)
}

func (o *funcWorkerOpt) Apply(opts *workerOpts) {
	o.f(opts)
}

func newWorkerOpt(f func(*workerOpts)) WorkerOpt {
	return &funcWorkerOpt{
		f: f,
	}
}

func WithDuration(duration time.Duration) WorkerOpt {
	return newWorkerOpt(func(opts *workerOpts) {
		opts.duration = duration
	})
}

func WithQpsWorker() WorkerOpt {
	return newWorkerOpt(func(opts *workerOpts) {
		opts.mode = WorkerModeQps
	})
}

func WithConstanceBucket(rate int) WorkerOpt {
	return newWorkerOpt(func(opts *workerOpts) {
		opts.bucket = BucketModeConstant
		opts.rate = rate
	})
}

func NewWorkerFactory() *workerFactory {
	return &workerFactory{}
}

func (f *workerFactory) newBucket(opts *workerOpts) (bucket.Bucket, error) {
	switch opts.bucket {
	case BucketModeConstant:
		return bucket.NewConstantBucket(opts.rate, opts.duration), nil
	case BucketModeUp:
		return nil, nil
	case BucketModelDown:
		return nil, nil
	case BucketModeRange:
		return nil, nil
	default:
		return nil, nil
	}
	return nil, nil
}

func (f *workerFactory) NewWorker(task task.Task, taskArgs interface{}, opts ...WorkerOpt) (worker Worker, err error) {
	workerOpts := defaultOpts()
	for _, opt := range opts {
		opt.Apply(workerOpts)
	}
	bucket, err := f.newBucket(workerOpts)
	if err != nil {
		return
	}
	switch workerOpts.mode {
	case WorkerModeQps:
		worker = newQpsWorker(bucket, task, taskArgs)
		return
	case WorkerModeConcurrency:
		return nil, &workerModeNotSupportErr{}
	default:
		return nil, &workerModeNotSupportErr{}
	}
	return
}
