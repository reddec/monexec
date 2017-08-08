package main

import (
	"time"
	"github.com/reddec/container"
	"github.com/reddec/container/plugin"
	"log"
	"os"
	"github.com/hashicorp/consul/api"
	"context"
	"sync"
	"github.com/Pallinder/go-randomdata"
	"errors"
	"io/ioutil"
	"strings"
	"gopkg.in/yaml.v2"
	"github.com/reddec/monexec/monexec"
	"path"
)

type Config struct {
	Services []monexec.Executable                       `yaml:"services"`
	Critical []string                                   `yaml:"critical,omitempty"`
	Consul struct {
		URL                       string        `yaml:"url"`
		TTL                       time.Duration `yaml:"ttl"`
		AutoDeregistrationTimeout time.Duration `yaml:"timeout"`
		Dynamic                   []string      `yaml:"register,omitempty"`
		Permanent                 []string      `yaml:"permanent,omitempty"`
	}                                           `yaml:"consul"`
	Telegram *Telegram                          `yaml:"telegram,omitempty"`
}

func (c *Config) MergeFrom(other *Config) error {
	c.Services = append(c.Services, other.Services...)
	c.Critical = append(c.Critical, other.Critical...)

	def := DefaultConfig()

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

	merged, err := mergeTelegram(c.Telegram, other.Telegram)
	if err != nil {
		return err
	}
	c.Telegram = merged
	return nil
}

func DefaultConfig() Config {
	config := Config{}
	// Default params for Consul
	config.Consul.TTL = 2 * time.Minute
	config.Consul.AutoDeregistrationTimeout = 5 * time.Minute
	config.Consul.URL = "http://localhost:8500"
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

	// -- telegram
	if config.Telegram != nil {
		err := config.Telegram.Prepare()
		if err != nil {
			log.Println("telegram plugin ont initialized due to", err)
		} else {
			sv.Events().AddHandler(config.Telegram)
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
			files = []os.FileInfo{stat}
		}
		for _, info := range files {
			if strings.HasSuffix(info.Name(), ".yml") || strings.HasSuffix(info.Name(), ".yaml") {
				data, err := ioutil.ReadFile(path.Join(location, info.Name()))
				if err != nil {
					return nil, err
				}
				var conf Config = DefaultConfig()
				err = yaml.Unmarshal(data, &conf)
				if err != nil {
					return nil, err
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
