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
	name  string
	path  string
	tasks []string
}

type PluginConfig struct {
	Name  string `json:"Name"`
	Path  string
	Tasks []string `json:"Tasks"`
}

func NewPlugin(config *PluginConfig) Plugin {
	return &plugin{
		name:  config.Name,
		path:  config.Path,
		tasks: config.Tasks,
	}
}

// TODO: thread-safe
func (p *plugin) Tasks() (map[string]task.Task, error) {
	client := server.NewPluginClient("unix:/tmp/gob/test.echo.socket")
	tasks := make(map[string]task.Task)
	for _, name := range p.tasks {
		handle := func(c context.Context, args string) error {
			// TODO: register on_start, on_success... events for collecting metrics
			return client.Call(c, name, args)
		}
		tasks[name] = task.NewTask(name, handle)
	}
	return tasks, nil
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
