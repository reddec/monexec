//

// +build !linux

package pool

import (
	"os/exec"
	"log"
)

func setAttrs(cmd *exec.Cmd) {
}

func kill(cmd *exec.Cmd, logger *log.Logger) error {
	return cmd.Process.Kill()
}
