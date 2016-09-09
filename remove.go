package vvmn

import (
	"fmt"
	"os"
	"path/filepath"
)

// Remove removes the specified Vim versions.
func Remove(versions ...string) error {
	for _, version := range versions {
		dir := filepath.Join(vvmnrootVim, version)
		srcDir := filepath.Join(vvmnrootSrc, version)
		if !exist(dir) {
			return fmt.Errorf("no vim version specified")
		}
		if err := os.RemoveAll(dir); err != nil {
			return err
		}
		if err := os.RemoveAll(srcDir); err != nil {
			return err
		}
	}
	return nil
}
