// Package vvmn provides support for managing Vim versions.
//
// First of all, SetRoot to determine root directory for vvmn:
//
//     vvmn.SetRoot("/path/to/root")
//
// Then, Get a specific Vim version:
//
//     err := vvmn.Get("vim1.7")
//     ...
//
package vvmn

import (
	"os"
	"path/filepath"
)

// Version is version string of vvmn.
const Version = "0.0.0"

// RepoURL indicates the Vim original repository url.
var RepoURL = "git://github.com/vim/vim.git"

var (
	vvmnroot     string
	vvmnrootSrc  string
	vvmnrootVim  string
	vvmnrootRepo string
)

// SetRoot sets root as vvmn's root directory.
func SetRoot(root string) {
	vvmnroot = root
	vvmnrootSrc = filepath.Join(root, "src")
	vvmnrootVim = filepath.Join(root, "vim")
	vvmnrootRepo = filepath.Join(root, "repo")
}

func exist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
