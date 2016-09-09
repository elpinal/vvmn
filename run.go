package vvmn

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Run executes the specified Vim version.
func Run(version string, args ...string) error {
	vimCmd := filepath.Join(vvmnrootVim, version, "bin", "vim")
	if !exist(vimCmd) {
		return fmt.Errorf("no installed vim version specified")
	}
	cmd := exec.Command(vimCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		if _, isExitError := err.(*exec.ExitError); !isExitError {
			return err
		}
	}
	return nil
}
