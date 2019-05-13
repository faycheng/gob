package plugin

import "github.com/faycheng/gob/task"

type PluginWatcher interface {
	Watch() (string, task.Task, error)
}

type pluginWatcher struct {
	root string
}

func NewPluginWatcher(root string) PluginWatcher {
	return &pluginWatcher{
		root: root,
	}
}

func (w *pluginWatcher) Watch() (key string, task task.Task, err error) {

	return
}
