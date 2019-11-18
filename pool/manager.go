package pool

import (
	"context"
	"sync"
)

type Manager struct {
	lock            sync.Mutex
	processManagers map[string]*ProcessManager
	plugins         []GlobalPlugin
}

func (mg *Manager) Plugin(plugin GlobalPlugin) *Manager {
	mg.plugins = append(mg.plugins, plugin)
	return mg
}

func (mg *Manager) Clean() {
	mg.lock.Lock()
	defer mg.lock.Unlock()
	// stop all
	wg := sync.WaitGroup{}
	for _, pm := range mg.processManagers {
		wg.Add(1)
		go func(pm *ProcessManager) {
			defer wg.Done()
			pm.Stop()
			for _, plugin := range mg.pluginsCopy() {
				plugin.AfterRemove(mg, pm)
			}
		}(pm)
	}
	wg.Wait()
	mg.processManagers = nil
}

func (mg *Manager) Add(ctx context.Context, template ExecutableTemplate) {
	mg.lock.Lock()
	defer mg.lock.Unlock()
	if mg.processManagers == nil {
		mg.processManagers = make(map[string]*ProcessManager)
	}

	if pm, exists := mg.processManagers[template.Name]; exists {
		pm.Stop()
		for _, plugin := range mg.pluginsCopy() {
			plugin.AfterRemove(mg, pm)
		}
	}
	pm := runProcessWithPreflightHandler(ctx, template, func(pm *ProcessManager) {
		for _, plugin := range mg.pluginsCopy() {
			plugin.BeforeAdd(ctx, mg, pm)
		}
	})
	mg.processManagers[template.Name] = pm
	for _, plugin := range mg.pluginsCopy() {
		plugin.AfterAdd(ctx, mg, pm)
	}
}

func (mg *Manager) Remove(name string) {
	mg.lock.Lock()
	defer mg.lock.Unlock()
	if process, exists := mg.processManagers[name]; exists {
		process.Stop()
		for _, plugin := range mg.pluginsCopy() {
			plugin.AfterRemove(mg, process)
		}
	}
	delete(mg.processManagers, name)
}

func (mg *Manager) Find(name string) *ProcessManager {
	mg.lock.Lock()
	defer mg.lock.Unlock()
	return mg.processManagers[name]
}

func (mg *Manager) pluginsCopy() []GlobalPlugin {
	return mg.plugins
}
