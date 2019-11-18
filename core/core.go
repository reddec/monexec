package core

import (
	"os/exec"
	"sync"
)

type Status int8

const (
	StatusStopped Status = iota
	StatusStarted
	StatusCrashed
)

type (
	OnStopped func(instance Instance)
	OnStarted func(instance Instance)
	OnCrashed func(instance Instance, err error)
)

type Instance interface {
	Start()         // start instance
	Stop()          // stop instance
	Restart()       // bounce instance
	Status() Status // current instance status
}

type instance struct {
	status Status
	lock   sync.RWMutex
	cmd    *exec.Cmd
}
