package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/susp/vvmn"
)

var cmdRun = &Command{
	Run:       runRun,
	UsageLine: "run version command [arguments]",
	Short:     "execute the specified Vim version",
	Long:      `Run executes the specified Vim version.`,
}

func runRun(cmd *Command, args []string) int {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", cmd.UsageLine)
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(cmd.Long))
		return 2
	}

	var vimArgs []string
	if len(args) != 1 {
		vimArgs = args[1:]
	}
	if err := vvmn.Run(args[0], vimArgs...); err != nil {
		log.Print(err)
		return 1
	}
	return 0
}
