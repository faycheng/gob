package task

import (
	"context"
	"fmt"

	"github.com/faycheng/goblinker"
)

type Signal string

const (
	OnStart   Signal = "on_start"
	OnSuccess        = "on_success"
	OnFailure        = "on_failure"
	OnError          = "on_error"
)

type Task interface {
	Connect(Signal, goblinker.Receiver) error
	Call(c context.Context, args interface{}) error
}

type task struct {
	key       string
	onStart   *goblinker.Signal
	onSuccess *goblinker.Signal
	onFailure *goblinker.Signal
	onError   *goblinker.Signal
}

// TODO: add handler
func NewTask(key string) Task {
	return &task{
		key:       key,
		onStart:   goblinker.NewSignal(true),
		onSuccess: goblinker.NewSignal(true),
		onFailure: goblinker.NewSignal(true),
		onError:   goblinker.NewSignal(true),
	}
}

func (t *task) Connect(signal Signal, receiver goblinker.Receiver) error {
	switch signal {
	case OnStart:
		t.onStart.Connect(receiver, "")
	case OnSuccess:
		t.onSuccess.Connect(receiver, "")
	case OnFailure:
		t.onFailure.Connect(receiver, "")
	case OnError:
		t.onError.Connect(receiver, "")
	default:
		return nil
	}
	return nil
}

func (t *task) Call(c context.Context, args interface{}) error {

	fmt.Printf("Task.Call(%+v)\n", args)
	return nil
}
