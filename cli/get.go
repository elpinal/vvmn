package cli

import (
	"log"

	"github.com/susp/vvmn"
)

var cmdGet = &Command{
	Run:       runGet,
	UsageLine: "get [-d] versions...",
	Short:     "download and install Vim",
	Long: `
Get downloads the specified Vim versions, and then installs them.

The -d flag instructs get to stop after downloading the Vim versions; that is,
it instructs get not to install the Vim versions.
`,
}

var (
	getD bool
)

func init() {
	cmdGet.Flag.BoolVar(&getD, "d", false, "")
}

// runGet executes get command and return exit code.
func runGet(cmd *Command, args []string) int {
	if len(args) == 0 {
		log.Print("vvmn get: no vim versions specified")
		return 1
	}

	for i, version := range args {
		if version == "latest" {
			latest, err := vvmn.LatestTag()
			if err != nil {
				log.Print(err)
				return 1
			}
			version = latest
			args[i] = latest
		}

		if err := vvmn.Download(version); err != nil {
			log.Print(err)
			return 1
		}
	}

	if getD {
		return 0
	}

	for _, version := range args {
		if err := vvmn.Install(version); err != nil {
			log.Print(err)
			return 1
		}
	}

	if err := vvmn.Use(args[len(args)-1:][0]); err != nil {
		log.Print(err)
		return 1
	}

	return 0
}
