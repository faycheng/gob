package plugin

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPlugin_Run(t *testing.T) {
	config := &PluginConfig{
		Name:  "echo",
		Path:  "../example/plugin/echo",
		Tasks: []string{"echo"},
	}
	plugin := NewPlugin(config)
	err := plugin.Run()
	assert.Empty(t, err)
	tasks, err := plugin.Tasks()
	assert.Empty(t, err)
	// Wait the sub-process of plugin
	time.Sleep(1 * time.Second)
	err = tasks["echo"].Call(context.TODO(), `{"msg": "hello world"}`)
	assert.Empty(t, err)
}
