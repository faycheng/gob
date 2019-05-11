package plugin

import "github.com/faycheng/gob/task"

type PluginWatcher interface {
	Watch() (task.Key, task.Task, error)
}

type pluginWatcher struct {
	root string
}

func NewPluginWatcher(root string) PluginWatcher {
	return &pluginWatcher{
		root: root,
	}
}

func (w *pluginWatcher) Watch() (key task.Key, task task.Task, err error) {

	return
}
