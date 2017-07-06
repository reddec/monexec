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

const (
	RestartInfinity       = -1              // Restart forever
	DefaultRestartTimeout = 5 * time.Second // Default timeout between restarts
	DefaultStopTimeout    = 5 * time.Second // Default timeout to wait for graceful shutdown
	DefaultStartTimeout   = 3 * time.Second // Default timeout to check process become alive
)

// Executable Event type
type EventType int

const (
	STARTED EventType = iota // Process started (after StartTimeout)
	STOPPED                  // Process stopped (no matter with error or not)
)

// Event of executable state
type Event struct {
	Type       EventType
	Executable *Executable
	Error      error
}

// Executable - basic information about process
type Executable struct {
	Label          string            `yaml:"label,omitempty"`           // Human-readable label for process. If not set - command used
	Command        string            `yaml:"command"`                   // Executable
	Args           []string          `yaml:"args,omitempty"`            // Arguments to command
	Environment    map[string]string `yaml:"environment,omitempty"`     // Additional environment variables
	Retries        int               `yaml:"retries,omitempty"`         // Restart retries limit. Negative value means infinity
	Critical       bool              `yaml:"critical,omitempty"`        // Stop all other processes on finish or error
	WorkDir        string            `yaml:"workdir,omitempty"`         // Working directory. If not set - current dir used
	StopTimeout    time.Duration     `yaml:"stop_timeout,omitempty"`    // Timeout before terminate process
	StartTimeout   time.Duration     `yaml:"start_timeout,omitempty"`   // Timeout to check process is still alive
	RestartTimeout time.Duration     `yaml:"restart_timeout,omitempty"` // Timeout before restart
}

// ID of process. By default Label is used, but if it not set, command name is selected
func (b *Executable) ID() string {
	id := b.Label
	if id == "" {
		id = b.Command
	}
	return id
}

// Mark process with custom label
func (b *Executable) Mark(label string) *Executable {
	b.Label = label
	return b
}

// Arg adds additional positional argument
func (b *Executable) Arg(arg string) *Executable {
	b.Args = append(b.Args, arg)
	return b
}

// Env adds additional environment key-value pair
func (b *Executable) Env(arg, value string) *Executable {
	if b.Environment == nil {
		b.Environment = make(map[string]string)
	}
	b.Environment[arg] = value
	return b
}

// try to do graceful process termination by sending SIGKILL. If no response after StopTimeout
// SIGTERM is used
func (exe *Executable) stopOrKill(logger *log.Logger, cmd *exec.Cmd) {
	ch := make(chan struct{}, 1)
	logger.Println("Sending SIGKTERM")
	cmd.Process.Signal(syscall.SIGTERM)
	go func() {
		cmd.Wait()
		ch <- struct{}{}
	}()

	select {
	case <-ch:
		logger.Println("Process gracefull stopped")
	case <-time.After(exe.StopTimeout):
		logger.Println("Process gracefull shutdown waiting timeout")
		cmd.Process.Signal(syscall.SIGKILL)
	}
}

// run once executable, wrap output and wait for finish
func (exe *Executable) runOnce(logger *log.Logger, stop <-chan struct{}) error {
	cmd := exec.Command(exe.Command, exe.Args...)
	for _, param := range os.Environ() {
		cmd.Env = append(cmd.Env, param)
	}
	if exe.Environment != nil {
		for k, v := range exe.Environment {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}
	if exe.WorkDir != "" {
		cmd.Dir = exe.WorkDir
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
	dumpDone := make(chan struct{}, 1)
	go func() {
		dumpToLogger(io.MultiReader(stdout, stderr), logger)
		dumpDone <- struct{}{}
	}()
	go func() { res <- cmd.Wait() }()
	select {
	case <-stop:
		exe.stopOrKill(logger, cmd)
		<-dumpDone
		return nil
	case err := <-res:
		<-dumpDone
		return err
	}
}

// Run loop that will monitor and restart process if needed
func (exe *Executable) Run(stop <-chan struct{}, events chan<- Event, logsink io.Writer) error {
	var err error
	if logsink == nil {
		logsink = os.Stdout
	}
	retries := exe.Retries
	logger := log.New(logsink, "["+exe.ID()+"] ", log.LstdFlags)

LOOP:
	for {
		started := make(chan error, 1)

		go func() { started <- exe.runOnce(logger, stop) }()

		select {
		case serr := <-started:
			err = serr
		case <-time.After(exe.StartTimeout):
			logger.Println("Process is still running for enough time - resetting restart limit")
			retries = exe.Retries
			events <- Event{STARTED, exe, nil}
			serr := <-started
			err = serr
		}
		events <- Event{STOPPED, exe, err}

		if err == nil {
			logger.Println("Process finished successfully")
		} else {
			logger.Println("Process finished unsuccessfully:", err)
		}
		if exe.Critical && err != nil {
			logger.Println("Process is critical - finishing work on error immediatly")
			break LOOP
		}

		if retries > 0 {
			logger.Println("Restarts left:", retries)
			retries--
		} else if retries == 0 {
			if exe.Retries != 0 { // not one-shot
				logger.Println("Process restart limit exceeded")
				if err == nil {
					err = errors.New("restart limit exceeded")
				}
			}
			break
		}

		logger.Println("Waiting", exe.RestartTimeout, "before restart")
		select {
		case <-stop:
			logger.Println("Gracefull shutdown restart loop")
			break LOOP
		case <-time.After(exe.RestartTimeout):
		}
	}
	if err != nil {
		logger.Println("Restart loop finished with error:", err)
	} else {
		logger.Println("Restart loop finished without error")
	}
	return err
}

// line-by-line writer to logger
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

// Monitor pool of executables
type Monitor struct {
	Executables []*Executable
	Logsink     io.Writer
}

// Add prepared executable
func (m *Monitor) Add(ex *Executable) *Executable {
	m.Executables = append(m.Executables, ex)
	return ex
}

// Oneshot adds an instance of non-critical non-restartable executable.
// Runs once with error tolerance
func (m *Monitor) Oneshot(command string, args ...string) *Executable {
	exe := &Executable{Retries: 0, Critical: false}
	exe.Command = command
	exe.Args = append(exe.Args, args...)
	exe.StopTimeout = DefaultStopTimeout
	return m.Add(exe)
}

// Critical adds an instance of non-restartable executable which will stop all another processes on error
func (m *Monitor) Critical(command string, args ...string) *Executable {
	exe := &Executable{Retries: 0, Critical: true}
	exe.Command = command
	exe.Args = append(exe.Args, args...)
	exe.StopTimeout = DefaultStopTimeout
	return m.Add(exe)
}

// Restart adds an instance of restartable non-critical executable
func (m *Monitor) Restart(maxRetries int, command string, args ...string) *Executable {
	exe := &Executable{Retries: maxRetries, Critical: false}
	exe.Command = command
	exe.Args = append(exe.Args, args...)
	exe.StartTimeout = DefaultStartTimeout
	exe.StopTimeout = DefaultStopTimeout
	exe.RestartTimeout = DefaultRestartTimeout
	return m.Add(exe)
}

// Forever adds an infinity-restartable non-critical executable
func (m *Monitor) Forever(command string, args ...string) *Executable {
	return m.Restart(RestartInfinity, command, args...)
}

// Run all processes and monitor them
func (m *Monitor) Run(ctx context.Context, events chan<- Event) error {
	done := make(chan struct{})

	wg := sync.WaitGroup{}
	wg.Add(len(m.Executables))
	errs := make([]error, len(m.Executables))
	for i, exe := range m.Executables {
		go func(i int, exe *Executable) {
			defer wg.Done()
			errs[i] = exe.Run(firstDone(ctx.Done(), done), events, m.Logsink)
			if errs[i] != nil && exe.Critical {
				close(done)
			}
		}(i, exe)
	}
	wg.Wait()
	return m.joinErrors(errs)
}

// RunNoEvent start Run with null events consumer
func (m *Monitor) RunNoEvents(ctx context.Context) error {
	dummy := make(chan Event, 1)
	go func() {
		for range dummy {
		}
	}()
	return m.Run(ctx, dummy)
}

// Start in background.
// Last event always will be with (nil) Executable, STOPPED state and error as result of Run
func (m *Monitor) Start(ctx context.Context) <-chan Event {
	ch := make(chan Event, 2*len(m.Executables))
	go func() {
		defer close(ch)
		err := m.Run(ctx, ch)
		ch <- Event{STOPPED, nil, err}
	}()
	return ch
}

func (m *Monitor) joinErrors(errs []error) error {
	errT := ""
	for i, err := range errs {
		if errT != "" {
			errT += "\n"
		}
		if err != nil {
			errT += m.Executables[i].ID() + ": " + err.Error()
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
