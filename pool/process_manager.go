package pool

import (
	"context"
	"time"
)

type ProcessStatus uint8

const (
	StatusScheduled = iota
	StatusRunning
	StatusRestarting
	StatusStopped
)

func RunProcess(ctx context.Context, template ExecutableTemplate) *ProcessManager {
	return runProcessWithPreflightHandler(ctx, template, func(pm *ProcessManager) {})
}

func runProcessWithPreflightHandler(ctx context.Context, template ExecutableTemplate, handlerFunc func(pm *ProcessManager)) *ProcessManager {
	pm := &ProcessManager{
		stop:     make(chan struct{}, 1),
		restart:  make(chan struct{}, 1),
		done:     make(chan struct{}),
		template: template,
	}
	handlerFunc(pm)
	go pm.loop(ctx)
	return pm
}

type ProcessManager struct {
	stop     chan struct{}
	done     chan struct{}
	restart  chan struct{}
	plugins  []ProcessPlugin
	template ExecutableTemplate
	status   ProcessStatus
}

func (pm *ProcessManager) Plugin(plugin ProcessPlugin) *ProcessManager {
	pm.plugins = append(pm.plugins, plugin)
	return pm
}

func (pm *ProcessManager) Template() ExecutableTemplate { return pm.template }

func (pm *ProcessManager) Status() ProcessStatus { return pm.status }

func (pm *ProcessManager) Stop() {
	select {
	case pm.stop <- struct{}{}:
	default:
	}
	<-pm.done
}

func (pm *ProcessManager) Restart() {
	select {
	case pm.restart <- struct{}{}:
	default:
	}
}

func (pm *ProcessManager) pluginsCopy() []ProcessPlugin { return pm.plugins }

func (pm *ProcessManager) loop(ctx context.Context) {
	defer func() {
		close(pm.done)
	}()
	pm.status = StatusScheduled
	for _, plugin := range pm.pluginsCopy() {
		plugin.Scheduled(ctx, pm)
	}
LOOP:
	for i := 0; pm.template.Restart == RestartAlways || i <= pm.template.Restart; i++ {
		process, err := pm.template.Start(ctx)
		if err != nil {
			pm.template.logger().Println("Failed to start:", err)
			goto WaitRestart
		}
		pm.status = StatusRunning
		for _, plugin := range pm.pluginsCopy() {
			plugin.Started(ctx, pm)
		}
		select {
		case <-process.Done():
		case <-pm.restart:
			process.Stop()
		case <-pm.stop:
			process.Stop()
			break LOOP
		case <-ctx.Done():
			process.Stop()
			break LOOP
		}

	WaitRestart:
		pm.status = StatusRestarting
		for _, plugin := range pm.pluginsCopy() {
			plugin.Restarting(ctx, pm)
		}
		select {
		case <-time.After(pm.template.RestartTimeout):
		case <-pm.stop:
			break LOOP
		case <-ctx.Done():
			break LOOP
		}
	}
	pm.status = StatusStopped
	for _, plugin := range pm.pluginsCopy() {
		plugin.Stopped(ctx, pm)
	}
}
