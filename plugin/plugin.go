package plugin

import (
	"context"
	"github.com/faycheng/gob/plugin/server"
	"github.com/faycheng/gob/task"
)

type Plugin interface {
	Tasks() (map[string]task.Task, error)
}

type plugin struct {
	name  string
	path  string
	tasks map[string]task.Task
}

type PluginConfig struct {
	name  string
	path  string
	tasks []string
}

func NewPlugin(config *PluginConfig) Plugin {
	client := server.NewPluginClient(config.path)
	tasks := make(map[string]task.Task)
	for _, name := range config.tasks {
		handle := func(c context.Context, args string) error {
			// TODO: register task events
			return client.Call(c, name, args)
		}
		tasks[name] = task.NewTask(name, handle)
	}
	return &plugin{
		name:  config.name,
		path:  config.path,
		tasks: tasks,
	}
}

// TODO: thread-safe
func (p *plugin) Tasks() (map[string]task.Task, error) {
	return p.tasks, nil
}
