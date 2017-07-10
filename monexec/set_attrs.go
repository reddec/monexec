// +build !linux
package monexec

import (
	"os/exec"
	"syscall"
)

func setAttrs(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}
