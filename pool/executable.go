package pool

import (
	"time"
	"os/exec"
	"os"
	"io"
	"log"
	"context"
	"strings"
	"path/filepath"
	"sync"
)

// Executable - basic information about process
type Executable struct {
	Name           string            `yaml:"label,omitempty"`         // Human-readable label for process. If not set - command used
	Command        string            `yaml:"command"`                 // Executable
	Args           []string          `yaml:"args,omitempty"`          // Arguments to command
	Environment    map[string]string `yaml:"environment,omitempty"`   // Additional environment variables
	WorkDir        string            `yaml:"workdir,omitempty"`       // Working directory. If not set - current dir used
	StopTimeout    time.Duration     `yaml:"stop_timeout,omitempty"`  // Timeout before terminate process
	RestartTimeout time.Duration     `yaml:"restart_delay,omitempty"` // Restart delay
	Restart        int               `yaml:"restart,omitempty"`       // How much restart allowed. -1 infinite
	LogFile        string            `yaml:"logFile,omitempty"`       // if empty - only to log. If not absolute - relative to workdir

	log        *log.Logger
	loggerInit sync.Once
}

func (b *Executable) WithName(name string) *Executable {
	cp := *b
	cp.loggerInit = sync.Once{}
	cp.Name = name
	return &cp
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

func (e *Executable) logger() *log.Logger {
	e.loggerInit.Do(func() {
		e.log = log.New(os.Stderr, "["+e.Name+"] ", log.LstdFlags)
	})
	return e.log
}

// try to do graceful process termination by sending SIGKILL. If no response after StopTimeout
// SIGTERM is used
func (exe *Executable) stopOrKill(cmd *exec.Cmd, res <-chan error) error {
	exe.logger().Println("Sending SIGINT")
	err := cmd.Process.Signal(os.Interrupt)
	if err != nil {
		exe.logger().Println("Failed send SIGINT:", err)
	}

	select {
	case err = <-res:
		exe.logger().Println("Process gracefull stopped")
	case <-time.After(exe.StopTimeout):
		exe.logger().Println("Process gracefull shutdown waiting timeout")
		err = kill(cmd, exe.logger())
	}
	return err
}

// run once executable, wrap output and wait for finish
func (exe *Executable) run(ctx context.Context) error {
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

	setAttrs(cmd)

	var outputs []io.Writer

	outputs = append(outputs, NewLoggerStream(exe.logger(), "out:"))

	res := make(chan error, 1)

	if exe.LogFile != "" {
		pth, _ := filepath.Abs(exe.LogFile)
		if pth != exe.LogFile {
			// relative
			wd, _ := filepath.Abs(exe.WorkDir)
			exe.LogFile = filepath.Join(wd, exe.LogFile)
		}
		logFile, err := os.OpenFile(exe.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			exe.logger().Println("Failed open log file:", err)
		} else {
			defer logFile.Close()
			outputs = append(outputs, logFile)
		}
	}

	logStream := io.MultiWriter(outputs...)

	cmd.Stderr = logStream
	cmd.Stdout = logStream

	err := cmd.Start()
	if err == nil {
		exe.logger().Println("Started with PID", cmd.Process.Pid)
	} else {
		exe.logger().Println("Failed start `", exe.Command, strings.Join(exe.Args, " "), "` :", err)
	}

	go func() { res <- cmd.Wait() }()
	select {
	case <-ctx.Done():
		err = exe.stopOrKill(cmd, res)
	case err = <-res:
	}
	return err
}

type runnable struct {
	Executable *Executable `json:"config"`
	Running    bool        `json:"running"`
	pool       *Pool
	closer     func()
	done       chan struct{}
}

func (exe *Executable) Start(ctx context.Context, pool *Pool) Instance {
	chCtx, closer := context.WithCancel(ctx)
	run := &runnable{
		Executable: exe,
		closer:     closer,
		done:       make(chan struct{}),
		pool:       pool,
	}
	go run.run(chCtx)
	return run
}

func (exe *Executable) Config() *Executable { return exe }

func (rn *runnable) run(ctx context.Context) {
	defer rn.closer()
	defer close(rn.done)
	restarts := rn.Executable.Restart
	rn.pool.OnSpawned(ctx, rn)
LOOP:
	for {
		rn.Running = true
		rn.pool.OnStarted(ctx, rn)
		err := rn.Executable.run(ctx)
		if err != nil {
			rn.Executable.logger().Println("stopped with error:", err)
		} else {
			rn.Executable.logger().Println("stopped")
		}
		rn.Running = false
		rn.pool.OnStopped(ctx, rn, err)
		if restarts != -1 {
			if restarts <= 0 {
				rn.Executable.logger().Println("max restarts attempts reached")
				break
			} else {
				restarts--
			}
		}
		rn.Executable.logger().Println("waiting", rn.Executable.RestartTimeout)
		select {
		case <-time.After(rn.Executable.RestartTimeout):
		case <-ctx.Done():
			rn.Executable.logger().Println("instance done:", ctx.Err())
			break LOOP
		}
	}
	rn.Executable.logger().Println("instance restart loop done")
	rn.pool.OnFinished(ctx, rn)
}

func (rn *runnable) Supervisor() Supervisor { return rn.Executable }

func (rn *runnable) Config() *Executable { return rn.Executable }

func (rn *runnable) Pool() *Pool { return rn.pool }

func (rn *runnable) Stop() {
	rn.closer()
	<-rn.done
}
