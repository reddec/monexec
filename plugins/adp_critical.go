package plugins

import (
	"context"
	"github.com/reddec/monexec/pool"
)

type Critical struct {
	Labels []string `mapstructure:"<ITEMS>"`
}

func (p *Critical) OnSpawned(ctx context.Context, sv pool.Instance) {}

func (p *Critical) OnStarted(ctx context.Context, sv pool.Instance) {}

func (p *Critical) OnStopped(ctx context.Context, sv pool.Instance, err error) {}

func (p *Critical) OnFinished(ctx context.Context, sv pool.Instance) {
	terminate := false
	for _, l := range p.Labels {
		if l == sv.Supervisor().Config().Name {
			terminate = true
			break
		}
	}
	if terminate {
		go sv.Pool().Terminate()
	}
}

func (a *Critical) MergeFrom(other interface{}) (error) {
	b := other.(*Critical)
	a.Labels = append(a.Labels, b.Labels...)
	return nil
}

func (a *Critical) Prepare(ctx context.Context, pl *pool.Pool) error {
	return nil
}
func (a *Critical) Close() error { return nil }
func init() {
	registerPlugin("critical", func(file string) PluginConfigNG {
		return &Critical{}
	})
}
