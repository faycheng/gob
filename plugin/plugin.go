package plugin

import (
	"context"
	"github.com/faycheng/gob/plugin/server"
	"github.com/faycheng/gob/task"
	"github.com/sirupsen/logrus"
	"os/exec"
)

type Plugin interface {
	Tasks() (map[string]task.Task, error)
	Run() error
}

type plugin struct {
	name       string
	path       string
	tasks      map[string]task.Task
	entrypoint string
}

type PluginConfig struct {
	name       string `json:"name"`
	path       string
	entrypoint string   `json:"entrypoint"`
	tasks      []string `json:"tasks"`
}

func NewPlugin(config *PluginConfig) Plugin {
	client := server.NewPluginClient(config.path)
	tasks := make(map[string]task.Task)
	for _, name := range config.tasks {
		handle := func(c context.Context, args string) error {
			// TODO: register on_start, on_success... events for collecting metrics
			return client.Call(c, name, args)
		}
		tasks[name] = task.NewTask(name, handle)
	}
	return &plugin{
		name:       config.name,
		path:       config.path,
		tasks:      tasks,
		entrypoint: config.entrypoint,
	}
}

// TODO: thread-safe
func (p *plugin) Tasks() (map[string]task.Task, error) {
	return p.tasks, nil
}

func (p *plugin) Run() error {
	cmd := exec.Command(p.entrypoint)
	if err := cmd.Start(); err != nil {
		return err
	}
	go func() {
		err := cmd.Wait()
		if err != nil {
			logrus.Errorf("plugin(%s) exit with err: %+v", p.entrypoint, err)
			return
		}
		logrus.Warnf("plugin(%s) exit", p.entrypoint)
	}()
	return nil
}
