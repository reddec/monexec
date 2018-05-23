package monexec

import (
	"time"
	"os/exec"
	"os"
	"io"
	"log"
	"bufio"
	"context"
	"github.com/reddec/container"
	"strings"
	"path/filepath"
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

// Factory of executables
func (e *Executable) Factory() (container.Runnable, error) {
	return &runnable{Executable: *e, logger: log.New(os.Stderr, "["+e.Name+"] ", log.LstdFlags)}, nil
}

type runnable struct {
	Executable
	logger *log.Logger
}

// ID of process. By default Label is used, but if it not set, command name is selected
func (b *runnable) Label() string {
	id := b.Name
	if id == "" {
		id = b.Command
	}
	return id
}

// try to do graceful process termination by sending SIGKILL. If no response after StopTimeout
// SIGTERM is used
func (exe *runnable) stopOrKill(logger *log.Logger, cmd *exec.Cmd, res <-chan error) error {
	logger.Println("Sending SIGINT")
	err := cmd.Process.Signal(os.Interrupt)
	if err != nil {
		logger.Println("Failed send SIGINT:", err)
	}

	select {
	case err = <-res:
		logger.Println("Process gracefull stopped")
	case <-time.After(exe.StopTimeout):
		logger.Println("Process gracefull shutdown waiting timeout")
		err = kill(cmd, logger)
	}
	return err
}

// run once executable, wrap output and wait for finish
func (exe *runnable) Run(ctx context.Context) error {
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

	outputs = append(outputs, NewLoggerStream(exe.logger, "out:"))

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
			exe.logger.Println("Failed open log file:", err)
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
		exe.logger.Println("Started with PID", cmd.Process.Pid)
	} else {
		exe.logger.Println("Failed start `", exe.Command, strings.Join(exe.Args, " "), "` :", err)
	}

	go func() { res <- cmd.Wait() }()
	select {
	case <-ctx.Done():
		err = exe.stopOrKill(exe.logger, cmd, res)
	case err = <-res:
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
