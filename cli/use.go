package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/susp/vvmn"
)

var cmdUse = &Command{
	Run:       runUse,
	UsageLine: "use version",
	Short:     "select a Vim version to use",
	Long:      `Use selects a Vim version to use.`,
}

// runUse executes use command and return exit code.
func runUse(cmd *Command, args []string) int {
	if len(args) == 0 {
		log.Print("vvmn use: no vim version specified")
		return 1
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", cmd.UsageLine)
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(cmd.Long))
		return 2
	}
	if err := vvmn.Use(args[0]); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
