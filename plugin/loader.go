package plugin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type PluginLoader interface {
	Load() (Plugin, error)
}

type pluginLoader struct {
	path string
}

func NewPluginLoader(path string) PluginLoader {
	return &pluginLoader{
		path: path,
	}
}

// TODO: wrap error
func (l *pluginLoader) Load() (plugin Plugin, err error) {
	// load plugin config
	fd, err := os.Open(fmt.Sprintf("%s/plugin.json", l.path))
	if err != nil {
		return
	}
	defer fd.Close()
	content, err := ioutil.ReadAll(fd)
	if err != nil {
		return
	}
	config := new(PluginConfig)
	err = json.Unmarshal(content, config)
	if err != nil {
		return
	}
	config.Path = l.path
	plugin = NewPlugin(config)
	return
}
