package task

import "sync"

type Key string

type Registry interface {
	Register(Key, Task) error
	Get(Key) (Task, bool)
}
type taskRegistry struct {
	taskSet sync.Map
}

func NewTaskRegistry() Registry {
	return &taskRegistry{}
}

func (r *taskRegistry) Register(key Key, task Task) error {
	r.taskSet.Store(key, task)
	return nil
}

func (r *taskRegistry) Get(key Key) (task Task, ok bool) {
	taskIface, ok := r.taskSet.Load(key)
	if !ok {
		return nil, false
	}
	task, ok = taskIface.(Task)
	return
}
