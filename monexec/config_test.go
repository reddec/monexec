package monexec

import (
	"testing"
	"github.com/reddec/container"
	"context"
	"time"
	"log"
	"os"
)

func TestConfig_Run(t *testing.T) {
	sv := container.NewSupervisor(log.New(os.Stderr, "[supervisor] ", log.LstdFlags))
	ctx, stop := context.WithCancel(context.Background())
	config := DefaultConfig()

	config.Services = append(config.Services, Executable{
		Name:    "srv1",
		Command: "echo",
		Args:    []string{"1234"},
	})
	config.Critical = append(config.Critical, "srv1")
	done := make(chan error, 1)
	go func() {
		done <- config.Run(sv, ctx)
	}()
	var err error
	select {
	case <-time.After(30 * time.Second):
		stop()
		err = <-done
	case err = <-done:
	}

	if err != nil {
		t.Error(err)
	}
}
