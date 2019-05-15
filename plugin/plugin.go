package plugin

import (
	"context"
	"fmt"
	"github.com/faycheng/gob/plugin/server"
	"github.com/faycheng/gob/task"
	"github.com/sirupsen/logrus"
	"os"
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
	entrypoint := fmt.Sprintf("%s/entrypoint", p.path)
	_, err := exec.LookPath(entrypoint)
	if err != nil {
		return err
	}
	cmd := exec.Command(entrypoint)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	go func() {
		err := cmd.Run()
		if err != nil {
			logrus.Errorf("plugin(%s/entrypoint.sh) exit with err: %+v", p.path, err)
			return
		}
		logrus.Infof("plugin(%s,entrypoint.sh) exit", p.path)
	}()
	return nil
}
