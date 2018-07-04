package monexec

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Pallinder/go-randomdata"
	"gopkg.in/yaml.v2"
	"github.com/reddec/monexec/pool"
	"github.com/mitchellh/mapstructure"
	"errors"
	"github.com/reddec/monexec/plugins"
	"reflect"
)

type Config struct {
	Services      []pool.Executable                 `yaml:"services"`
	Plugins       map[string]interface{}            `yaml:",inline"` // all unparsed means plugins
	loadedPlugins map[string]plugins.PluginConfigNG `yaml:"-"`
}

func (c *Config) MergeFrom(other *Config) error {
	c.Services = append(c.Services, other.Services...)
	// -- merge plugins
	for otherPluginName, otherPluginInstance := range other.loadedPlugins {
		if ownPlugin, needMerge := c.loadedPlugins[otherPluginName]; needMerge {
			err := ownPlugin.MergeFrom(otherPluginInstance)
			if err != nil {
				return errors.New("merge " + otherPluginName + ": " + err.Error())
			}
		} else { // new one - just copy
			c.loadedPlugins[otherPluginName] = otherPluginInstance
		}
	}
	return nil
}

func (c *Config) ClosePlugins() {
	for _, plugin := range c.loadedPlugins {
		plugin.Close()
	}
}

func DefaultConfig() Config {
	config := Config{}

	config.loadedPlugins = make(map[string]plugins.PluginConfigNG)
	return config
}

func FillDefaultExecutable(exec *pool.Executable) {
	if exec.RestartTimeout == 0 {
		exec.RestartTimeout = 6 * time.Second
	}
	if exec.Restart == 0 {
		exec.Restart = -1
	}
	if exec.StopTimeout == 0 {
		exec.StopTimeout = 3 * time.Second
	}
	if exec.Name == "" {
		exec.Name = randomdata.Noun() + "-" + randomdata.Adjective()
	}
}

func (config *Config) Run(sv *pool.Pool, ctx context.Context) error {
	// Initialize plugins
	//-- prepare and add all plugins
	for pluginName, pluginInstance := range config.loadedPlugins {
		err := pluginInstance.Prepare(ctx, sv)
		if err != nil {
			log.Println("failed prepare plugin", pluginName, "-", err)
		} else {
			log.Println("plugin", pluginName, "ready")
			sv.Watch(pluginInstance)
		}
	}

	// Run
	for i := range config.Services {
		exec := config.Services[i]
		FillDefaultExecutable(&exec)
		sv.Add(&exec)
	}

	sv.StartAll(ctx)
	return nil
}

func LoadConfig(locations ...string) (*Config, error) {
	c := DefaultConfig()
	ans := &c
	var files []os.FileInfo
	for _, location := range locations {
		if stat, err := os.Stat(location); err != nil {
			return nil, err
		} else if stat.IsDir() {
			fs, err := ioutil.ReadDir(location)
			if err != nil {
				return nil, err
			}
			files = fs
		} else {
			location = filepath.Dir(location)
			files = []os.FileInfo{stat}
		}
		for _, info := range files {
			if strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml") {
				fileName := path.Join(location, info.Name())
				data, err := ioutil.ReadFile(fileName)
				if err != nil {
					return nil, err
				}
				var conf = DefaultConfig()
				err = yaml.Unmarshal(data, &conf)
				if err != nil {
					return nil, err
				}

				// -- load all plugins for current config here
				for pluginName, description := range conf.Plugins {
					pluginInstance, found := plugins.BuildPlugin(pluginName, fileName)
					if !found {
						log.Println("plugin", pluginName, "not found")
						continue
					}

					var wrap = description
					refVal := reflect.ValueOf(wrap)
					if wrap != nil && refVal.Type().Kind() == reflect.Slice {
						wrap = map[string]interface{}{
							"<ITEMS>": description,
						}
					}

					config := &mapstructure.DecoderConfig{
						Metadata:   nil,
						Result:     pluginInstance,
						DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
					}

					decoder, err := mapstructure.NewDecoder(config)
					if err != nil {
						panic(err) // failed to initialize decoder - something really wrong
					}

					err = decoder.Decode(wrap)
					if err != nil {
						log.Println("failed load plugin", pluginName, "-", err)
						continue
					}
					conf.loadedPlugins[pluginName] = pluginInstance
				}

				err = ans.MergeFrom(&conf)
				if err != nil {
					return nil, err
				}
			}
		}
	}
	return ans, nil
}
