package monexec

import (
	"time"
	"os/exec"
	"syscall"
	"os"
	"io"
	"log"
	"bufio"
	"sync"
	"context"
	"errors"
)

const RestartInfinity = -1
const DefaultRestartTimeout = 5 * time.Second
const DefaultStopTimeout = 5 * time.Second

type BaseInfo struct {
	ID          string            `yaml:"id"`
	Command     string            `yaml:"command"`
	Args        []string          `yaml:"args"`
	Environment map[string]string `yaml:"environment"`
	Workdir     string            `yaml:"workdir"`
	StopTimeout time.Duration     `yaml:"stop_timeout"`
}

func (b *BaseInfo) GetID() string {
	id := b.ID
	if id == "" {
		id = b.Command
	}
	return id
}

func (b *BaseInfo) WithID(id string) *BaseInfo {
	b.ID = id
	return b
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

type Delayed struct {
	BaseInfo                   `yaml:",squash"`
	RestartDelay time.Duration `yaml:"restart_delay"`
}

type Executable struct {
	Delayed       `yaml:",squash"`
	Retries  int  `yaml:"retries"`
	Critical bool `yaml:"critical"`
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
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGTERM,
	}
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

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}
	defer cmd.Process.Kill()
	logger.Println("Started with PID", cmd.Process.Pid)
	res := make(chan error, 1)
	go dumpToLogger(io.MultiReader(stdout, stderr), logger)
	go func() { res <- cmd.Wait() }()
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
	logger := log.New(os.Stdout, "["+exe.GetID()+"] ", log.LstdFlags)
LOOP:
	for {
		err = exe.Run(logger, stop)
		if err == nil {
			logger.Println("Process finished successfully")
		} else {
			logger.Println("Process finished unsuccessfully:", err)
		}
		if exe.Critical && err != nil {
			logger.Println("Process is critical - finishing work on error immediatly")
			break LOOP
		}

		if exe.Retries > 0 {
			logger.Println("Restarts left:", exe.Retries)
			exe.Retries--
		} else if exe.Retries == 0 {
			logger.Println("Process restart limit exceeded")
			break
		}

		logger.Println("Waiting", exe.RestartDelay, "before restart")
		select {
		case <-stop:
			logger.Println("Gracefull shutdown restart loop")
			break LOOP
		case <-time.After(exe.RestartDelay):
		}
	}
	if err != nil {
		logger.Println("Restart loop finished with error:", err)
	} else {
		logger.Println("Restart loop finished without error")
	}
	return err
}

func dumpToLogger(reader io.Reader, logger *log.Logger) {
	scanner := bufio.NewReader(reader)
	for {
		line, _, err := scanner.ReadLine()
		if err != nil {
			break
		}
		logger.Println("out:", string(line))
	}
}

type Monitor struct {
	Executables []*Executable
}

func (m *Monitor) Add(ex *Executable) *Executable {
	m.Executables = append(m.Executables, ex)
	return ex
}

func (m *Monitor) Oneshot(command string, args ...string) *Executable {
	exe := &Executable{Retries: 0, Critical: false}
	exe.Command = command
	exe.Args = append(exe.Args, args...)
	exe.StopTimeout = DefaultStopTimeout
	return m.Add(exe)
}

func (m *Monitor) Critical(command string, args ...string) *Executable {
	exe := &Executable{Retries: 0, Critical: true}
	exe.Command = command
	exe.Args = append(exe.Args, args...)
	exe.StopTimeout = DefaultStopTimeout
	return m.Add(exe)
}

func (m *Monitor) Restart(maxRetries int, command string, args ...string) *Executable {
	exe := &Executable{Retries: maxRetries, Critical: false}
	exe.Command = command
	exe.Args = append(exe.Args, args...)
	exe.StopTimeout = DefaultStopTimeout
	exe.RestartDelay = DefaultRestartTimeout
	return m.Add(exe)
}

func (m *Monitor) Forever(command string, args ...string) *Executable {
	return m.Restart(RestartInfinity, command, args...)
}

func (m *Monitor) Run(ctx context.Context) error {
	done := make(chan struct{})

	wg := sync.WaitGroup{}
	errs := make([]error, len(m.Executables))
	for i, exe := range m.Executables {
		wg.Add(1)
		go func(i int, exe *Executable) {
			defer wg.Done()
			errs[i] = exe.Daemonize(firstDone(ctx.Done(), done))
			if errs[i] != nil && exe.Critical {
				close(done)
			}
		}(i, exe)
	}
	wg.Wait()
	return m.joinErrors(errs)
}

func (m *Monitor) joinErrors(errs []error) error {
	errT := ""
	for i, err := range errs {
		if errT != "" {
			errT += "\n"
		}
		if err != nil {
			errT += m.Executables[i].GetID() + ": " + err.Error()
		}
	}
	if errT != "" {
		return errors.New(errT)
	}
	return nil
}

func firstDone(chanA, chanB <-chan struct{}) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		select {
		case <-chanA:
			close(ch)
		case <-chanB:
			close(ch)
		}
	}()
	return ch
}
