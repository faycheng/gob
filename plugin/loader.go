package plugin

import (
	"context"
	"github.com/faycheng/gob/task"
)

type PluginLoader interface {
	Load() (Plugin, error)
}

type pluginLoaer struct {
	path string
}

func NewPluginLoader(path string) PluginLoader {
	return &pluginLoaer{
		path: path,
	}
}

// TODO: load plugin
func (l *pluginLoaer) Load() (plugin Plugin, err error) {
	taskSet = make(map[string]task.Task, 0)
	taskSet["echo"] = task.NewTask("echo")
	// TODO: register on_start, on_success... events for collecting metrics
	taskSet["echo"].Connect(task.OnStart, func(c context.Context, args ...interface{}) error {
		return nil
	})

	return
}
