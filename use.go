package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var cmdUse = &Command{
	Run:       runUse,
	UsageLine: "use ",
	Short:     "Use Vim",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdUse.Flag.BoolVar(&flagA, "a", false, "")
}

// runUse executes use command and return exit code.
func runUse(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "vvmn use: no vim version specified")
		return 1
	}
	currentDir := filepath.Join(VvmnDir, "vims", "current")
	version := args[0]
	if version == "system" {
		if _, err := os.Stat(currentDir); err != nil {
			return 0
		}
		if err := os.RemoveAll(currentDir); err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed use system version of vim"))
			return 1
		}
		return 0
	}
	vimsDir := filepath.Join(VvmnDir, "vims", version)
	if _, err := os.Stat(vimsDir); err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "no installed version of vim specified"))
		return 1
	}
	if _, err := os.Stat(currentDir); err == nil {
		if err := os.RemoveAll(currentDir); err != nil {
			fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed unuse former version of vim"))
			return 1
		}
	}
	err := os.Symlink(vimsDir, currentDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed create symbolic link"))
		return 1
	}

	return 0
}
