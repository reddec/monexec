# Sample plugin

```go
package plugins

import (
	"context"
	"github.com/reddec/monexec/pool"
)


type MyPlugin struct {}

func (p *MyPlugin) Prepare(ctx context.Context, pl *pool.Pool) error { return nil }

func (p *MyPlugin) OnSpawned(ctx context.Context, sv pool.Instance) {}

func (p *MyPlugin) OnStarted(ctx context.Context, sv pool.Instance) {}

func (p *MyPlugin) OnStopped(ctx context.Context, sv pool.Instance, err error) {}

func (p *MyPlugin) OnFinished(ctx context.Context, sv pool.Instance) {}

func (p *MyPlugin) MergeFrom(other interface{}) error { return nil}

func (a *MyPlugin) Close() error { return nil }

func init() {
    registerPlugin("myPlugin", func(file string) PluginConfigNG {
        return &MyPlugin{}
    })
}
```