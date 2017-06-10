package cli

import (
	"log"
	"strings"

	"github.com/elpinal/vvmn"
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
	i := 0
	for i < len(args) && !strings.HasPrefix(args[i], "-") {
		i++
	}
	versions, cmdArgs := args[:i], args[i:]
	if len(versions) == 0 {
		log.Print("vvmn get: no vim versions specified")
		return 1
	}
	for n, version := range versions {
		if version == "latest" {
			latest, err := vvmn.LatestTag()
			if err != nil {
				log.Print(err)
				return 1
			}
			versions[n] = latest
		}
	}

	for _, version := range versions {
		if err := vvmn.Download(version); err != nil {
			log.Print(err)
			return 1
		}
	}

	if getD {
		return 0
	}

	for _, version := range versions {
		if err := vvmn.Install(version, cmdArgs...); err != nil {
			log.Print(err)
			return 1
		}
	}

	if err := vvmn.Use(versions[len(versions)-1:][0]); err != nil {
		log.Print(err)
		return 1
	}

	return 0
}
