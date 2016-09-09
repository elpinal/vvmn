package vvmn

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Use selects a Vim version to use.
func Use(version string) error {
	currentDir := filepath.Join(vvmnrootVim, "current")
	versionsDir := filepath.Join(vvmnrootVim, version)
	if !exist(versionsDir) {
		return fmt.Errorf("no installed vim version specified")
	}
	if err := os.RemoveAll(currentDir); err != nil {
		return errors.Wrap(err, "failed to stop using former vim version")
	}
	err := os.Symlink(versionsDir, currentDir)
	if err != nil {
		return errors.Wrap(err, "failed to create symbolic link")
	}
	return nil
}
