package plugin

import "github.com/faycheng/gob/task"

type PluginLoader interface {
	Load() (map[task.Key]task.Task, error)
}

type pluginLoaer struct {
	path string
}

func NewPluginLoader(path string) PluginLoader {
	return &pluginLoaer{
		path: path,
	}
}

func (l *pluginLoaer) Load() (taskSet map[task.Key]task.Task, err error) {

	return
}
