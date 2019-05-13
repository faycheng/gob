package plugin

import (
	"github.com/faycheng/gob/task"
	"google.golang.org/grpc"
)

type Plugin interface {
	Tasks() (map[string]task.Task, error)
}

type plugin struct {
	// NOTE: grpc client
	client *grpc.ClientConn
	task   map[string]task.Task
}

func NewPlugin(socket string) Plugin {
	return &plugin{}
}

// TODO: parse tasks
func (p *plugin) Tasks() (map[string]task.Task, error) {
	// client.Tasks()
	// new task
	// connect events
	return nil, nil
}
