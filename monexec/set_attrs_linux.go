package monexec

import (
	"syscall"
	"os/exec"
	"log"
)

func setAttrs(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Pdeathsig: syscall.SIGKILL,
		Setpgid:   true,
	}
}

func kill(cmd *exec.Cmd, logger *log.Logger) error {
	pgid, err := syscall.Getpgid(cmd.Process.Pid)
	if err == nil {
		if err := syscall.Kill(-pgid, syscall.SIGKILL); err != nil {
			logger.Println("Failed kill by process group:", err)
			err = cmd.Process.Kill() // fallback
		}
	} else {
		err = cmd.Process.Kill() // fallback
	}

	if err != nil {
		logger.Println("Failed kill:", err)
	}
	return err
}
