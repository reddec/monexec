package monexec

import (
	"time"
	"os/exec"
	"syscall"
	"os"
	"io"
	"log"
	"bufio"
	"fmt"
	"sync"
	"context"
	"errors"
)

const RestartInfinity = -1

type BaseInfo struct {
	ID          string
	Command     string
	Args        []string
	Environment map[string]string
	Workdir     string
	StopTimeout time.Duration
}

type Delayed struct {
	BaseInfo
	RestartDelay time.Duration
}

func (b *BaseInfo) Arg(arg string) *BaseInfo {
	b.Args = append(b.Args, arg)
	return b
}

func (b *BaseInfo) Env(arg, value string) *BaseInfo {
	if b.Environment == nil {
		b.Environment = make(map[string]string)
	}
	b.Environment[arg] = value
	return b
}

type Executable struct {
	Delayed
	RestartCount int
	Critical     bool
}

func (exe *Executable) stopOrKill(logger *log.Logger, cmd *exec.Cmd) {
	ch := make(chan struct{}, 1)
	logger.Println("Sending SIGKILL")
	cmd.Process.Kill()
	go func() {
		cmd.Wait()
		ch <- struct{}{}
	}()

	select {
	case <-ch:
		logger.Println("Process gracefull stopped")
	case <-time.After(exe.StopTimeout):
		logger.Println("Process gracefull shutdown waiting timeout")
		cmd.Process.Signal(syscall.SIGTERM)
	}
}

func (exe *Executable) Run(logger *log.Logger, stop <-chan struct{}) error {
	cmd := exec.Command(exe.Command, exe.Args...)
	cmd.SysProcAttr.Pdeathsig = syscall.SIGTERM
	for _, param := range os.Environ() {
		cmd.Env = append(cmd.Env, param)
	}
	if exe.Environment != nil {
		for k, v := range exe.Environment {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}
	if exe.Workdir != "" {
		cmd.Dir = exe.Workdir
	}
	err := cmd.Start()
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()
	logger.SetPrefix(logger.Prefix() + fmt.Sprintf("[%v] ", cmd.Process.Pid))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	res := make(chan error, 1)
	go dumpToLogger(io.MultiReader(stdout, stderr), logger)
	go func() { res <- cmd.Run() }()
	select {
	case <-stop:
		exe.stopOrKill(logger, cmd)
		return nil
	case err := <-res:
		return err
	}
}

func (exe *Executable) Daemonize(stop <-chan struct{}) error {
	var err error
	logger := log.New(os.Stdout, "["+exe.ID+"] ", log.LstdFlags)
	for {
		err = exe.Run(logger, stop)
		if err == nil {
			logger.Println("Process finished successfully")
		} else {
			logger.Println("Process finished unsuccessfully:", err)
		}
		if exe.Critical && err != nil {
			logger.Println("Process is critical - finishing work on error immediatly")
			return err
		}

		if exe.RestartCount > 0 {
			logger.Println("Restarts left:", exe.RestartCount)
			exe.RestartCount--
		} else if exe.RestartCount == 0 {
			logger.Println("Process restart limit exceeded")
			break
		}

		logger.Println("Waiting", exe.RestartDelay, "before restart")
		select {
		case <-stop:
			log.Println("Gracefull shutdown restart loop")
			return err
		case <-time.After(exe.RestartDelay):
		}
	}
	return err
}

func dumpToLogger(reader io.Reader, logger *log.Logger) {
	scanner := bufio.NewReader(reader)
	for {
		line, _, err := scanner.ReadLine()
		logger.Println([]string(line))
		if err != nil {
			break
		}
	}
}

type Monitor struct {
	Executables []Executable
}

func (m *Monitor) Add(ex Executable) *Monitor {
	m.Executables = append(m.Executables, ex)
	return m
}

func (m *Monitor) Oneshot(base *BaseInfo) *Monitor {
	exe := Executable{RestartCount: 0, Critical: false}
	exe.BaseInfo = *base
	return m.Add(exe)
}

func (m *Monitor) Critical(base *BaseInfo) *Monitor {
	exe := Executable{RestartCount: 0, Critical: true}
	exe.BaseInfo = *base
	return m.Add(exe)
}

func (m *Monitor) Restart(base *Delayed, maxRetries int) *Monitor {
	exe := Executable{RestartCount: maxRetries, Critical: false}
	exe.Delayed = *base
	return m.Add(exe)
}

func (m *Monitor) Forever(base *Delayed) *Monitor {
	return m.Restart(base, RestartInfinity)
}

func (m *Monitor) Run(ctx context.Context) error {
	wg := sync.WaitGroup{}
	errs := make([]error, len(m.Executables))
	for i, exe := range m.Executables {
		wg.Add(1)
		go func(i int, exe *Executable) {
			defer wg.Done()
			errs[i] = exe.Daemonize(ctx.Done())
		}(i, &exe)
	}
	wg.Wait()
	return joinErrors(errs...)
}

func joinErrors(errs ...error) error {
	errT := ""
	for _, err := range errs {
		if errT != "" {
			errT += "\n"
		}
		if err != nil {
			errT += err.Error()
		}
	}
	if errT != "" {
		return errors.New(errT)
	}
	return nil
}
