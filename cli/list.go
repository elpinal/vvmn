package cli

import (
	"log"

	"github.com/susp/vvmn"
)

var cmdList = &Command{
	Run:       runList,
	UsageLine: "list",
	Short:     "list installed Vim versions",
	Long: `
List lists installed Vim versions.
Vim versions is divided into states.
	`,
}

func init() {
	// Set your flag here like below.
	// cmdList.Flag.BoolVar(&flagA, "a", false, "")
}

// runList executes list command and return exit code.
func runList(cmd *Command, args []string) int {
	if len(args) != 0 {
		cmd.Usage()
		return 2
	}

	logger := log.New(cmd.OutStream, "", 0)

	list := vvmn.List()

	if list.Current != "" {
		logger.Print("Current:")
		logger.Print()
		logger.Print("\t", list.Current)
		logger.Print()
	}

	if len(list.Installed) > 0 {
		logger.Print("Installed:")
		logger.Print()
		for _, installed := range list.Installed {
			logger.Print("\t", installed)
		}
		logger.Print()
	}

	if len(list.Downloaded) > 0 {
		logger.Print("Just downloaded; not installed:")
		logger.Print()
		for _, downloaded := range list.Downloaded {
			logger.Print("\t", downloaded)
		}
		logger.Print()
	}

	return 0
}
