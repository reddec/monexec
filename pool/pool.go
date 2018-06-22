package pool

import (
	"context"
	"sync"
)

type Instance interface {
	Stop()
	Config() *Executable
	Supervisor() Supervisor
	Pool() *Pool
}

type Supervisor interface {
	Start(ctx context.Context, pool *Pool) Instance
	Config() *Executable
}

type EventHandler interface {
	OnSpawned(ctx context.Context, in Instance)
	OnStarted(ctx context.Context, in Instance)
	OnStopped(ctx context.Context, in Instance, err error)
	OnFinished(ctx context.Context, in Instance)
}

type Pool struct {
	handlers     []EventHandler
	handlersLock sync.RWMutex

	supervisors []Supervisor
	svLock      sync.RWMutex

	instances []Instance
	inLock    sync.RWMutex

	doneInit sync.Once
	done     chan struct{}

	terminating bool
}

func (p *Pool) StopAll() {
	wg := sync.WaitGroup{}
	for _, sv := range p.grabInstances() {
		wg.Add(1)
		go func(sv Instance) {
			defer wg.Done()
			p.Stop(sv)
		}(sv)
	}
	wg.Wait()
}

func (p *Pool) StartAll(ctx context.Context) {
	if p.terminating {
		return
	}
	wg := sync.WaitGroup{}
	for _, sv := range p.Supervisors() {
		wg.Add(1)
		go func(sv Supervisor) {
			defer wg.Done()
			p.Start(ctx, sv)
		}(sv)
	}
	wg.Wait()
}

func (p *Pool) Start(ctx context.Context, sv Supervisor) Instance {
	if p.terminating {
		return nil
	}

	in := sv.Start(ctx, p)
	p.inLock.Lock()
	p.instances = append(p.instances, in)
	p.inLock.Unlock()
	return in
}

func (p *Pool) Stop(in Instance) {
	in.Stop()
	p.inLock.Lock()
	for i, v := range p.instances {
		if v == in {
			p.instances = append(p.instances[:i], p.instances[i+1:]...)
			break
		}
	}
	p.inLock.Unlock()
}

func (p *Pool) Add(sv Supervisor) {
	if p.terminating {
		return
	}
	p.svLock.Lock()
	defer p.svLock.Unlock()
	p.supervisors = append(p.supervisors, sv)
}

func (p *Pool) Watch(handler EventHandler) {
	p.handlersLock.Lock()
	defer p.handlersLock.Unlock()
	p.handlers = append(p.handlers, handler)
}

func (p *Pool) cloneHandlers() []EventHandler {
	p.handlersLock.RLock()
	var dest = make([]EventHandler, len(p.handlers))
	copy(dest, p.handlers)
	p.handlersLock.RUnlock()
	return dest
}

func (p *Pool) Supervisors() []Supervisor {
	p.svLock.RLock()
	var dest = make([]Supervisor, len(p.supervisors))
	copy(dest, p.supervisors)
	p.svLock.RUnlock()
	return dest
}

func (p *Pool) Instances() []Instance {
	p.inLock.RLock()
	var dest = make([]Instance, len(p.instances))
	copy(dest, p.instances)
	p.inLock.RUnlock()
	return dest
}

func (p *Pool) grabInstances() []Instance {
	p.inLock.Lock()
	var dest = p.instances
	p.instances = nil
	p.inLock.Unlock()
	return dest
}

func (p *Pool) OnSpawned(ctx context.Context, sv Instance) {
	for _, handler := range p.cloneHandlers() {
		handler.OnSpawned(ctx, sv)
	}
}

func (p *Pool) OnStarted(ctx context.Context, sv Instance) {
	for _, handler := range p.cloneHandlers() {
		handler.OnStarted(ctx, sv)
	}
}

func (p *Pool) OnStopped(ctx context.Context, sv Instance, err error) {
	for _, handler := range p.cloneHandlers() {
		handler.OnStopped(ctx, sv, err)
	}
}

func (p *Pool) OnFinished(ctx context.Context, sv Instance) {
	for _, handler := range p.cloneHandlers() {
		handler.OnFinished(ctx, sv)
	}
}

func (p *Pool) doneChan() chan struct{} {
	p.doneInit.Do(func() {
		p.done = make(chan struct{}, 1)
	})
	return p.done
}

func (p *Pool) notifyDone() {
	close(p.doneChan())
}

func (p *Pool) Done() <-chan struct{} {
	return p.doneChan()
}

func (p *Pool) Terminate() {
	if p.terminating {
		return
	}
	p.terminating = true
	p.StopAll()
	p.notifyDone()
}
