//

// +build !linux

package pool

import (
	"log"
	"os/exec"
)

func setAttrs(cmd *exec.Cmd) {
}

func kill(cmd *exec.Cmd, logger *log.Logger) {
	err := cmd.Process.Kill()
	if err != nil {
		logger.Println("Failed kill:", err)
	}
}
