package plugins

import (
	"fmt"
	"github.com/phuslu/goproxy/context"
	"net/http"
)

type Plugin interface {
	PluginName() string
	Fetch(*context.Context, *http.Request) (*http.Response, error)
}

type RegisteredPlugin struct {
	New func() (Plugin, error)
}

var (
	plugins map[string]*RegisteredPlugin
)

func init() {
	plugins = make(map[string]*RegisteredPlugin)
}

// Register a Plugin
func Register(name string, registeredPlugin *RegisteredPlugin) error {
	if _, exists := plugins[name]; exists {
		return fmt.Errorf("Name already registered %s", name)
	}

	plugins[name] = registeredPlugin
	return nil
}

// NewPlugin creates a new Plugin of type "name"
func NewPlugin(name string) (Plugin, error) {
	plugin, exists := plugins[name]
	if !exists {
		return nil, fmt.Errorf("hosts: Unknown plugin %q", name)
	}
	return plugin.New()
}
