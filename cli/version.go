package cli

import (
	"fmt"
	"os"
	"strings"

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
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", cmd.UsageLine)
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(cmd.Long))
		return 2
	}

	fmt.Printf("vvmn version %s\n", vvmn.Version)
	return 0
}
