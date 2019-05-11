package task

import (
	"context"
	"github.com/faycheng/goblinker"
)

type signal string

const (
	OnStart   signal = "on_start"
	OnSuccess        = "on_success"
	OnFailure        = "on_failure"
	OnError          = "on_error"
)

type Task interface {
	Connect(signal, goblinker.Receiver) error
	Call(c context.Context, args interface{}) error
}
