package plugins

import "github.com/reddec/container"

// Base interface for any future plugins
type PluginConfig interface {
	// Must be monitor
	container.Monitor
	// Merge change from other instance. Other is always has same type as original
	MergeFrom(other interface{}) error
	// Prepare internal state
	Prepare() error
}

// factories of plugins
var plugins = make(map[string]func(fileName string) PluginConfig)

// Register one plugin factory. File name not for parsing!
func registerPlugin(name string, factory func(fileName string) PluginConfig) {
	plugins[name] = factory
}

// Build but not fill one config
func BuildPlugin(name string, file string) (PluginConfig, bool) {
	if plugin, ok := plugins[name]; ok {
		return plugin(file), true
	}
	return nil, false
}
