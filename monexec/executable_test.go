package monexec

import (
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	"time"
)

func TestMonitor_Run(t *testing.T) {
	m := Monitor{}
	m.Oneshot("echo", "123", "456").Mark("test1")
	m.Restart(3, "nc", "-l", "0").Mark("netcat")
	m.Critical("/bin/sh", "-c", "sleep 2").Mark("shell-delay")
	ctx, stp := context.WithCancel(context.Background())
	go func() {
		<-time.After(5 * time.Second)
		stp()
	}()
	err := m.RunNoEvents(ctx)

	assert.NoError(t, err)
}
