package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

var cmdList = &Command{
	Run:       runList,
	UsageLine: "list ",
	Short:     "List Vim",
	Long: `

	`,
}

func init() {
	// Set your flag here like below.
	// cmdList.Flag.BoolVar(&flagA, "a", false, "")
}

// runList executes uninstall command and return exit code.
func runList(args []string) int {
	current, _ := os.Readlink(filepath.Join(VvmnDir, "vims", "current"))
	currentVersion := filepath.Base(current)
	vims, err := ioutil.ReadDir(filepath.Join(VvmnDir, "vims"))
	if err != nil {
		fmt.Fprintln(os.Stderr, errors.Wrap(err, "failed list versions of vim"))
	}
	for _, vim := range vims {
		version := vim.Name()
		if version == "current" {
			continue
		}
		var mark string
		if version == currentVersion {
			mark = "*"
		} else {
			mark = " "
		}
		fmt.Println(mark, version)
	}
	return 0
}
