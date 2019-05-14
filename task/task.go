package task

import (
	"context"
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
	handle    func(c context.Context, args string) error
	onStart   *goblinker.Signal
	onSuccess *goblinker.Signal
	onFailure *goblinker.Signal
	onError   *goblinker.Signal
}

func NewTask(key string, handle func(c context.Context, args string) error) Task {
	return &task{
		key:       key,
		handle:    handle,
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
	// TODO: send signal
	err := t.handle(c, args.(string))
	return err
}
