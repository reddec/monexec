package plugins

import (
	"github.com/reddec/monexec/pool"
	"io"
	"context"
)

// factories of plugins
var plugins = make(map[string]func(fileName string) PluginConfigNG)

// Register one plugin factory. File name not for parsing!
func registerPlugin(name string, factory func(fileName string) PluginConfigNG) {
	plugins[name] = factory
}

// Build but not fill one config
func BuildPlugin(name string, file string) (PluginConfigNG, bool) {
	if plugin, ok := plugins[name]; ok {
		return plugin(file), true
	}
	return nil, false
}

// Base interface for any future plugins
type PluginConfigNG interface {
	// Must handle events
	pool.EventHandler
	// Closable
	io.Closer
	// Merge change from other instance. Other is always has same type as original
	MergeFrom(other interface{}) error
	// Prepare internal state
	Prepare(ctx context.Context, pl *pool.Pool) error
}
