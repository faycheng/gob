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

type pluginLoaer struct {
	path string
}

func NewPluginLoader(path string) PluginLoader {
	return &pluginLoaer{
		path: path,
	}
}

// TODO: load plugin
// TODO: wrap error
func (l *pluginLoaer) Load() (plugin Plugin, err error) {
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
	// TODO: custom plugin path
	config.path = "/Users/chengfei/Dropbox/workspace/golang/src/github.com/faycheng/gob/example/plugin/echo"
	plugin = NewPlugin(config)
	return
}
