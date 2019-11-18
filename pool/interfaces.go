package pool

import (
	"context"
	"io"
)

type ProcessPlugin interface {
	io.Closer
	Scheduled(ctx context.Context, processManager *ProcessManager)
	Started(ctx context.Context, processManager *ProcessManager)
	Restarting(ctx context.Context, processManager *ProcessManager)
	Stopped(ctx context.Context, processManager *ProcessManager)
}

type GlobalPlugin interface {
	io.Closer
	BeforeAdd(ctx context.Context, manager *Manager, pm *ProcessManager)
	AfterAdd(ctx context.Context, manager *Manager, pm *ProcessManager)
	AfterRemove(manager *Manager, pm *ProcessManager)
}
