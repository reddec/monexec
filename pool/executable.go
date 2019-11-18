package pool

import (
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

func (template ExecutableTemplate) prepareCommand() *exec.Cmd {
	cmd := exec.Command(template.Command, template.Args...)
	for _, param := range os.Environ() {
		cmd.Env = append(cmd.Env, param)
	}
	if template.Environment != nil {
		for k, v := range template.Environment {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}
	for _, fileName := range template.EnvFiles {
		params, err := ParseEnvironmentFile(fileName)
		if err != nil {
			template.log.Println("failed parse environment file", fileName, ":", err)
			continue
		}
		for k, v := range params {
			cmd.Env = append(cmd.Env, k+"="+v)
		}
	}
	if template.WorkDir != "" {
		cmd.Dir = template.WorkDir
	}

	setAttrs(cmd)
	return cmd
}

func (template ExecutableTemplate) prepareOutputs(logger *log.Logger) (stdout, stderr []io.WriteCloser) {
	output := NewLoggerStream(logger, "out:")
	stderr = append(stderr, output)
	stdout = append(stdout, output)

	if template.LogFile != "" {
		logFile, err := os.OpenFile(template.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			template.log.Println("Failed open log file:", err)
		} else {
			stderr = append(stderr, logFile)
			stdout = append(stdout, logFile)
		}
	}
	return
}

func (template ExecutableTemplate) Start(ctx context.Context) (*Process, error) {
	proc := &Process{
		template: template,
		stopRQ:   make(chan struct{}, 1),
		done:     make(chan struct{}),
	}

	cmd := template.prepareCommand()
	stdout, stderr := template.prepareOutputs(template.logger())
	defer func() {
		for _, obj := range stdout {
			_ = obj.Close()
		}
		for _, obj := range stderr {
			_ = obj.Close()
		}
	}()

	logStderrStream := io.MultiWriter(mapWriterCloserToWriter(stderr)...)

	var logStdoutStream io.Writer
	if template.RawOutput {
		var output []io.Writer
		output = append(output, mapWriterCloserToWriter(stdout)...)
		output = append(output, os.Stdout)
		logStdoutStream = io.MultiWriter(output...)
	} else {
		logStdoutStream = io.MultiWriter(mapWriterCloserToWriter(stdout)...)
	}

	cmd.Stderr = logStderrStream
	cmd.Stdout = logStdoutStream

	err := cmd.Start()
	if err != nil {
		return nil, err
	}
	go proc.loop(ctx, cmd, template.logger())
	return proc, nil
}

type Process struct {
	template ExecutableTemplate
	stopRQ   chan struct{}
	done     chan struct{}
	err      error
}

func (proc *Process) Done() <-chan struct{}        { return proc.done }
func (proc *Process) Error() error                 { return proc.err }
func (proc *Process) Template() ExecutableTemplate { return proc.template }

func (proc *Process) Stop() {
	select {
	case proc.stopRQ <- struct{}{}:
	default:
	}
	<-proc.done
}

func (proc *Process) loop(ctx context.Context, cmd *exec.Cmd, logger *log.Logger) {
	stopped := make(chan error, 1)
	defer func() {
		close(proc.done)
	}()

	go func() {
		stopped <- cmd.Wait()
		close(stopped)
	}()

	var alreadyStopped bool
	select {
	case err := <-stopped: // executable stopped
		if err != nil {
			logger.Println("Stopped with error:", err)
		} else {
			logger.Println("Stopped without error")
		}
		alreadyStopped = true
	case <-proc.stopRQ: // requested to stop
	case <-ctx.Done():
	}
	var err error
	if !alreadyStopped {
		err = proc.gracefulStop(cmd, stopped, logger)
	}
	proc.err = err
}

func (proc *Process) gracefulStop(cmd *exec.Cmd, stopped chan error, logger *log.Logger) error {
	interrupt(cmd, logger)
	select {
	case <-time.After(proc.template.StopTimeout):
		kill(cmd, logger)
		return <-stopped
	case err := <-stopped:
		return err
	}
}

func interrupt(cmd *exec.Cmd, logger *log.Logger) {
	logger.Println("Sending SIGINT")
	err := cmd.Process.Signal(os.Interrupt)
	if err != nil {
		logger.Println("Failed send SIGINT:", err)
	}
}

func mapWriterCloserToWriter(closers []io.WriteCloser) []io.Writer {
	var ans = make([]io.Writer, len(closers))
	for i, v := range closers {
		ans[i] = v
	}
	return ans
}
