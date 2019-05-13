package worker

import (
	"context"
	"sync"

	"github.com/faycheng/gob/bucket"
	"github.com/faycheng/gob/task"
)

type Worker interface {
	Run()
}

type qpsWorker struct {
	stop     bool
	task     task.Task
	bucket   bucket.Bucket
	taskArgs interface{}
}

func newQpsWorker(bucket bucket.Bucket, task task.Task, taskArgs interface{}) Worker {
	return &qpsWorker{
		task:     task,
		bucket:   bucket,
		taskArgs: taskArgs,
	}
}

func (qw *qpsWorker) Run() {
	var wg sync.WaitGroup
	for qw.bucket.Get() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := qw.task.Call(context.TODO(), qw.taskArgs)
			if err != nil {
				panic(err)
			}
		}()
	}
	wg.Wait()
}
