package main

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Pallinder/go-randomdata"
	"github.com/hashicorp/consul/api"
	"github.com/reddec/container"
	"github.com/reddec/container/plugin"
	"github.com/reddec/monexec/monexec"
	"gopkg.in/yaml.v2"
	"github.com/reddec/monexec/plugins"
	"github.com/mitchellh/mapstructure"
	"errors"
)

type ConsulConfig struct {
	URL                       string        `yaml:"url"`
	TTL                       time.Duration `yaml:"ttl"`
	AutoDeregistrationTimeout time.Duration `yaml:"timeout"`
	Dynamic                   []string      `yaml:"register,omitempty"`
	Permanent                 []string      `yaml:"permanent,omitempty"`
}

func DefaultConsulConfig() *ConsulConfig {
	// Default params for Consul
	return &ConsulConfig{
		TTL:                       2 * time.Minute,
		AutoDeregistrationTimeout: 5 * time.Minute,
		URL:                       "http://localhost:8500",
	}
}

type Config struct {
	Services      []monexec.Executable            `yaml:"services"`
	Critical      []string                        `yaml:"critical,omitempty"`
	Consul        *ConsulConfig                   `yaml:"consul,omitempty"`
	Plugins       map[string]interface{}          `yaml:",inline"` // all unparsed means plugins
	loadedPlugins map[string]plugins.PluginConfig `yaml:"-"`
}

func (c *Config) MergeFrom(other *Config) error {
	c.Services = append(c.Services, other.Services...)
	c.Critical = append(c.Critical, other.Critical...)

	def := DefaultConfig()
	if c.Consul == nil {
		c.Consul = other.Consul
	} else if other.Consul != nil {
		if c.Consul.URL == def.Consul.URL {
			c.Consul.URL = other.Consul.URL
		} else if c.Consul.URL != def.Consul.URL && other.Consul.URL != def.Consul.URL && other.Consul.URL != c.Consul.URL {
			return errors.New("Different CONSUL definition (different URL) - specify same or only once")
		}

		if c.Consul.TTL == def.Consul.TTL {
			c.Consul.TTL = other.Consul.TTL
		} else if c.Consul.TTL != def.Consul.TTL && other.Consul.TTL != def.Consul.TTL && other.Consul.TTL != c.Consul.TTL {
			return errors.New("Different CONSUL definition (different TTL) - specify same or only once")
		}

		if c.Consul.AutoDeregistrationTimeout == def.Consul.AutoDeregistrationTimeout {
			c.Consul.AutoDeregistrationTimeout = other.Consul.AutoDeregistrationTimeout
		} else if c.Consul.AutoDeregistrationTimeout != def.Consul.AutoDeregistrationTimeout &&
			other.Consul.AutoDeregistrationTimeout != def.Consul.AutoDeregistrationTimeout &&
			other.Consul.AutoDeregistrationTimeout != c.Consul.AutoDeregistrationTimeout {
			return errors.New("Different CONSUL definition (different AutoDeregistrationTimeout) - specify same or only once")
		}

		c.Consul.Permanent = append(c.Consul.Permanent, other.Consul.Permanent...)
		c.Consul.Dynamic = append(c.Consul.Dynamic, other.Consul.Dynamic...)
	}
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

func DefaultConfig() Config {
	config := Config{}

	config.loadedPlugins = make(map[string]plugins.PluginConfig)
	return config
}

func FillDefaultExecutable(exec *monexec.Executable) {
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

func (config *Config) Run(sv container.Supervisor, ctx context.Context) error {
	critical := plugin.NewCritical(sv, log.New(os.Stderr, "[critical-plugin] ", log.LstdFlags), config.Critical...)
	sv.Events().AddHandler(critical)

	// Initialize plugins
	// -- consul
	if config.Consul != nil {
		consulConfig := api.DefaultConfig()
		consulConfig.Address = config.Consul.URL

		consul, err := api.NewClient(consulConfig)
		if err != nil {
			return err
		}
		var consulRegs []plugin.ConsulRegistration
		for _, label := range config.Consul.Dynamic {
			consulRegs = append(consulRegs, plugin.ConsulRegistration{Permanent: false, Label: label})
		}
		for _, label := range config.Consul.Permanent {
			consulRegs = append(consulRegs, plugin.ConsulRegistration{Permanent: true, Label: label})
		}
		consulLogger := log.New(os.Stderr, "[consul-plugin] ", log.LstdFlags)
		consulService := plugin.NewConsul(consul, config.Consul.TTL, config.Consul.AutoDeregistrationTimeout, consulLogger, consulRegs)
		defer consulService.Close()
		sv.Events().AddHandler(consulService)
	}

	//-- prepare and add all plugins
	for pluginName, pluginInstance := range config.loadedPlugins {
		err := pluginInstance.Prepare()
		if err != nil {
			log.Println("failed prepare plugin", pluginName, "-", err)
		} else {
			log.Println("plugin", pluginName, "ready")
			sv.Events().AddHandler(pluginInstance)
		}
	}

	// Run
	wg := sync.WaitGroup{}
	for _, exec := range config.Services {
		FillDefaultExecutable(&exec)
		wg.Add(1)
		go func(exec monexec.Executable) {
			defer wg.Done()
			container.Wait(sv.Watch(ctx, exec.Factory, exec.Restart, exec.RestartTimeout, false))
		}(exec)
	}

	wg.Wait()

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
					err := mapstructure.Decode(description, pluginInstance)
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
