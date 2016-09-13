package cli

import (
	"log"

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
		cmd.Usage()
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
