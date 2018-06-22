package plugins

import (
	"time"

	"github.com/hashicorp/consul/api"
	"log"
	"sync"
	"context"
	"fmt"
	"os"
	"errors"
	"github.com/reddec/monexec/pool"
)

type ConsulPlugin struct {
	URL                       string        `yaml:"url"`
	TTL                       time.Duration `yaml:"ttl"`
	AutoDeregistrationTimeout time.Duration `yaml:"timeout"`
	Dynamic                   []string      `yaml:"register,omitempty"`
	Permanent                 []string      `yaml:"permanent,omitempty"`

	registerLabels map[string]consulRegistration `yaml:"-"`
	Log            *log.Logger                   `yaml:"-"`
	Client         *api.Client                   `yaml:"-"`
	matched        map[string]struct{}           `yaml:"-"`
	lock           sync.Mutex                    `yaml:"-"`
	stop           chan struct{}                 `yaml:"-"`
	done           chan struct{}                 `yaml:"-"`
}

func DefaultConsul() ConsulPlugin {
	return ConsulPlugin{
		URL:                       "http://localhost:8500",
		AutoDeregistrationTimeout: 5 * time.Minute,
		TTL:                       2 * time.Minute,
	}
}

func (p *ConsulPlugin) Prepare(ctx context.Context, pl *pool.Pool) error {
	consulConfig := api.DefaultConfig()
	consulConfig.Address = p.URL
	consul, err := api.NewClient(consulConfig)
	if err != nil {
		return err
	}
	var consulRegs []consulRegistration
	for _, label := range p.Dynamic {
		consulRegs = append(consulRegs, consulRegistration{Permanent: false, Label: label})
	}
	for _, label := range p.Permanent {
		consulRegs = append(consulRegs, consulRegistration{Permanent: true, Label: label})
	}

	lbs := make(map[string]consulRegistration)
	for _, v := range consulRegs {
		lbs[v.Label] = v
	}

	p.Log = log.New(os.Stderr, "[consul] ", log.LstdFlags)
	p.stop = make(chan struct{}, 1)
	p.done = make(chan struct{}, 1)
	p.matched = make(map[string]struct{})
	p.registerLabels = lbs
	p.Client = consul

	go p.checkLoop()
	return nil
}

func (p *ConsulPlugin) OnSpawned(ctx context.Context, sv pool.Instance) {}

func (c *ConsulPlugin) OnStarted(ctx context.Context, sv pool.Instance) {
	label := sv.Config().Name
	info, exists := c.registerLabels[label]
	if !exists {
		return
	}
	dereg := c.AutoDeregistrationTimeout
	if dereg < c.TTL {
		dereg = 2 * c.TTL
	}
	if dereg < 1*time.Minute {
		dereg = 1 * time.Minute
	}
	err := c.Client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name: label,
		Tags: []string{fmt.Sprintf("%v", os.Getpid())},
		Check: &api.AgentServiceCheck{
			TTL:                            c.TTL.String(),
			DeregisterCriticalServiceAfter: dereg.String(),
		},
	})
	if err != nil {
		c.Log.Println("Can't register service", label, "in Consul:", err)
	} else {
		checkID := label + ":ttl"
		reg := api.AgentCheckRegistration{}
		reg.Name = checkID
		reg.TTL = c.TTL.String()
		reg.ServiceID = label

		if !info.Permanent {
			reg.DeregisterCriticalServiceAfter = dereg.String()
		}

		err = c.Client.Agent().CheckRegister(&reg)
		if err != nil {
			c.Log.Println("Can't register service TTL check", label, "in Consul:", err)
		} else {
			c.Log.Println("Service", label, "registered in Consul")
			c.lock.Lock()
			c.matched[checkID] = struct{}{}
			c.lock.Unlock()
		}
	}

}

func (c *ConsulPlugin) OnStopped(ctx context.Context, sv pool.Instance, err error) {
	label := sv.Config().Name
	info, exists := c.registerLabels[label]
	if !exists {
		return
	}
	c.lock.Lock()
	delete(c.matched, label+":ttl")
	c.lock.Unlock()

	if !info.Permanent {
		err = c.Client.Agent().ServiceDeregister(label)
		if err != nil {
			c.Log.Println("Can't deregister service", label, "in Consul:", err)
		} else {
			c.Log.Println("Service", label, "deregistered in Consul")
		}
	}
}

func (p *ConsulPlugin) OnFinished(ctx context.Context, sv pool.Instance) {}

func (p *ConsulPlugin) MergeFrom(a interface{}) error {
	other := a.(*ConsulPlugin)
	def := DefaultConsul()

	if p.URL == def.URL {
		p.URL = other.URL
	} else if p.URL != def.URL && other.URL != def.URL && other.URL != p.URL {
		return errors.New("Different CONSUL definition (different URL) - specify same or only once")
	}

	if p.TTL == def.TTL {
		p.TTL = other.TTL
	} else if p.TTL != def.TTL && other.TTL != def.TTL && other.TTL != p.TTL {
		return errors.New("Different CONSUL definition (different TTL) - specify same or only once")
	}

	if p.AutoDeregistrationTimeout == def.AutoDeregistrationTimeout {
		p.AutoDeregistrationTimeout = other.AutoDeregistrationTimeout
	} else if p.AutoDeregistrationTimeout != def.AutoDeregistrationTimeout &&
		other.AutoDeregistrationTimeout != def.AutoDeregistrationTimeout &&
		other.AutoDeregistrationTimeout != p.AutoDeregistrationTimeout {
		return errors.New("Different CONSUL definition (different AutoDeregistrationTimeout) - specify same or only once")
	}

	p.Permanent = append(p.Permanent, other.Permanent...)
	p.Dynamic = append(p.Dynamic, other.Dynamic...)
	return nil

}

func (c *ConsulPlugin) checkLoop() {
	defer close(c.done)
LOOP:
	for {
		select {
		case <-time.After(c.TTL / 2):
			c.updateChecks()
		case <-c.stop:
			break LOOP
		}
	}
}

func (c *ConsulPlugin) updateChecks() {
	c.lock.Lock()
	wg := sync.WaitGroup{}
	wg.Add(len(c.matched))
	for id, _ := range c.matched {
		go func(id string) {
			defer wg.Done()
			err := c.Client.Agent().UpdateTTL(id, "application running", "pass")
			if err != nil {
				c.Log.Println("Can't update TTL for service", id, "in Consul:", err)
			}
		}(id)
	}
	c.lock.Unlock()
	wg.Wait()
}

func (c *ConsulPlugin) Close() error {
	close(c.stop)
	<-c.done
	return nil
}

// Define options for consul registration
type consulRegistration struct {
	// Keep service registered even if stopped
	Permanent bool `json:"permanent,omitempty" yaml:"permanent,omitempty" ini:"permanent,omitempty"`
	// Name of service
	Label string `json:"label" yaml:"label" ini:"label"`
}

func init() {
	registerPlugin("consul", func(file string) PluginConfigNG {
		x := DefaultConsul()
		return &x
	})
}
