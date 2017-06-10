package cli

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mitchellh/go-homedir"
	"github.com/elpinal/vvmn"
)

type CLI struct {
	OutStream io.Writer
	ErrStream io.Writer
}

func (c CLI) Run(args []string) int {
	flags := flag.NewFlagSet("vvmn", flag.ContinueOnError)
	flags.Usage = c.usage
	flags.SetOutput(c.ErrStream)
	if err := flags.Parse(args); err != nil {
		return 2
	}
	log.SetFlags(0)
	log.SetOutput(c.ErrStream)

	args = flags.Args()
	if len(args) < 1 {
		c.usage()
		return 2
	}

	if args[0] == "help" {
		return c.help(args[1:])
	}

	if vvmn.GetRoot() == "" {
		if root := os.Getenv("VVMNROOT"); root != "" {
			vvmn.SetRoot(root)
		} else {
			home, err := homedir.Dir()
			if err != nil {
				log.Print(err)
				return 2
			}
			vvmn.SetRoot(filepath.Join(home, ".vvmn"))
		}
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.OutStream = c.OutStream
			cmd.ErrStream = c.ErrStream

			cmd.Flag.Init(args[0], flag.ContinueOnError)
			cmd.Flag.SetOutput(c.ErrStream)

			cmd.Flag.Usage = func() { cmd.Usage() }

			if err := cmd.Flag.Parse(args[1:]); err != nil {
				return 2
			}
			args = cmd.Flag.Args()

			return cmd.Run(cmd, args)
		}
	}

	log.Printf("vvmn: unknown subcommand %q\nRun 'vvmn help' for usage.\n", args[0])
	return 2
}

// A Command is an implementation of a vvmn command
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string) int

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'vvmn help' output.
	Short string

	// Long is the long message shown in the 'vvmn help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	OutStream io.Writer
	ErrStream io.Writer
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(c.ErrStream, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(c.ErrStream, "%s\n", strings.TrimSpace(c.Long))
}

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by 'vvmn help'.
var commands = []*Command{
	cmdGet,
	cmdList,
	cmdUse,
	cmdRemove,
	cmdRun,
	cmdVersion,
}

var usageTemplate = `vvmn is a tool for managing Vim versions.

Usage:

	vvmn command [arguments]

The commands are:
{{range .}}
	{{.Name | printf "%-11s"}} {{.Short}}{{end}}

Use "vvmn help [command]" for more information about a command.

`

var helpTemplate = `usage: vvmn {{.UsageLine}}

{{.Long | trim}}
`

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": strings.TrimSpace})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		panic(err)
	}
}

func printUsage(w io.Writer) {
	bw := bufio.NewWriter(w)
	tmpl(bw, usageTemplate, commands)
	bw.Flush()
}

func (c CLI) usage() {
	printUsage(c.ErrStream)
}

// help implements the 'help' command.
func (c CLI) help(args []string) int {
	if len(args) == 0 {
		printUsage(c.OutStream)
		// not exit 2: succeeded at 'vvmn help'.
		return 0
	}
	if len(args) != 1 {
		fmt.Fprintf(c.ErrStream, "usage: vvmn help command\n\nToo many arguments given.\n")
		return 2 // failed at 'vvmn help'
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(c.OutStream, helpTemplate, cmd)
			// not exit 2: succeeded at 'vvmn help cmd'.
			return 0
		}
	}

	fmt.Fprintf(c.ErrStream, "Unknown help topic %#q.  Run 'vvmn help'.\n", arg)
	return 2 // failed at 'vvmn help cmd'
}
