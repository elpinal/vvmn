package cli

import (
	"fmt"
	"log"
	"os"
	"strings"

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

func doOnce(f func()) func() {
	var done bool
	return func() {
		if done {
			return
		}
		f()
		done = true
	}
}

func genHeader(header string) func() {
	return func() {
		log.Print()
		log.Print(header)
		log.Print()
	}
}

// runList executes list command and return exit code.
func runList(cmd *Command, args []string) int {
	if len(args) != 0 {
		fmt.Fprintf(os.Stderr, "usage: %s\n\n", cmd.UsageLine)
		fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(cmd.Long))
		return 2
	}

	list := vvmn.List()
	if list == nil {
		return 0
	}

	for _, info := range list {
		if info.Current {
			log.Print("Current:")
			log.Print()
			log.Print("\t", info.Name)
			break
		}
	}

	ih := doOnce(genHeader("Installed:"))
	for _, info := range list {
		if !info.Current && info.Installed {
			ih()
			log.Print("\t", info.Name)
		}
	}

	dh := doOnce(genHeader("Just downloaded; not installed:"))
	for _, info := range list {
		if !info.Current && !info.Installed {
			dh()
			log.Print("\t", info.Name)
		}
	}

	log.Print()

	return 0
}
