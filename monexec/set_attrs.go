package monexec

import (
	"os/exec"
	"syscall"
)

// +build !windows
func setAttrs(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGKILL,
	}
}
