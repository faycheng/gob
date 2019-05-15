package plugin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPluginLoader_Load(t *testing.T) {
	loader := NewPluginLoader("../example/plugin/echo")
	_, err := loader.Load()
	assert.Equal(t, nil, err)
}
