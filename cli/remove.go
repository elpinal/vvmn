package cli

import (
	"log"

	"github.com/susp/vvmn"
)

var cmdRemove = &Command{
	Run:       runRemove,
	UsageLine: "remove versions...",
	Short:     "remove Vim versions",
	Long:      `Remove removes the specified Vim versions.`,
}

// runRemove executes remove command and return exit code.
func runRemove(cmd *Command, args []string) int {
	if len(args) == 0 {
		log.Print("vvmn remove: no vim versions specified")
		return 1
	}
	if err := vvmn.Remove(args...); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
