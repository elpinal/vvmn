package cli

import (
	"fmt"

	"github.com/susp/vvmn"
)

var cmdVersion = &Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "print vvmn version",
	Long:      `Version prints the vvmn version.`,
}

func runVersion(cmd *Command, args []string) int {
	if len(args) != 0 {
		cmd.Usage()
		return 2
	}

	fmt.Fprintf(cmd.OutStream, "vvmn version %s\n", vvmn.Version)
	return 0
}
