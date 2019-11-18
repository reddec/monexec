package pool

import (
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const RestartAlways = -1

// Executable - basic information about process.
type ExecutableTemplate struct {
	Name           string            `yaml:"label,omitempty"`         // Human-readable label for process. If not set - command used
	Command        string            `yaml:"command"`                 // Executable
	Args           []string          `yaml:"args,omitempty"`          // Arguments to command
	Environment    map[string]string `yaml:"environment,omitempty"`   // Additional environment variables
	EnvFiles       []string          `yaml:"envFiles"`                // Additional environment variables from files (not found files ignored). Format key=value
	WorkDir        string            `yaml:"workdir,omitempty"`       // Working directory. If not set - current dir used
	StopTimeout    time.Duration     `yaml:"stop_timeout,omitempty"`  // Timeout before terminate process
	RestartTimeout time.Duration     `yaml:"restart_delay,omitempty"` // Restart delay
	Restart        int               `yaml:"restart,omitempty"`       // How much restart allowed. -1 infinite
	LogFile        string            `yaml:"logFile,omitempty"`       // if empty - only to log. If not absolute - relative to workdir
	RawOutput      bool              `yaml:"raw,omitempty"`           // print stdout as-is without prefixes
	Watch          string            `yaml:"watch,omitempty"`         // watch file or directory for changes. If changed - service will be restarted

	log        *log.Logger
	loggerInit sync.Once
}

func (template *ExecutableTemplate) WithName(name string) *ExecutableTemplate {
	cp := *template
	cp.loggerInit = sync.Once{}
	cp.Name = name
	return &cp
}

// Arg adds additional positional argument
func (template *ExecutableTemplate) Arg(arg string) *ExecutableTemplate {
	template.Args = append(template.Args, arg)
	return template
}

// Env adds additional environment key-value pair
func (template *ExecutableTemplate) Env(arg, value string) *ExecutableTemplate {
	if template.Environment == nil {
		template.Environment = make(map[string]string)
	}
	template.Environment[arg] = value
	return template
}

func (template *ExecutableTemplate) logger() *log.Logger {
	template.loggerInit.Do(func() {
		template.resolve()
		template.log = log.New(os.Stderr, "["+template.Name+"] ", log.LstdFlags)
	})
	return template.log
}

func (template *ExecutableTemplate) resolve() {
	if template.LogFile != "" {
		pth, _ := filepath.Abs(template.LogFile)
		if pth != template.LogFile {
			// relative
			wd, _ := filepath.Abs(template.WorkDir)
			template.LogFile = filepath.Join(wd, template.LogFile)
		}
	}
}
