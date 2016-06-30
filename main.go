package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mitchellh/go-homedir"
)

// A Command is an implementation of a vvmn command
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(args []string) int

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'vvmn help' output.
	Short string

	// Long is the long message shown in the 'vvmn help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet
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
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

// Commands lists the available commands and help topics.
// The order here is the order in which they are printed by 'vvmn help'.
var commands = []*Command{
	cmdList,
	cmdInstall,
	cmdUse,
	cmdUninstall,
}

var RepoURL = "git://github.com/vim/vim.git"
var VvmnDir string

func main() {

	flag.Usage = usage
	flag.Parse()
	log.SetFlags(0)

	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}

	home, err := homedir.Dir()
	VvmnDir = filepath.Join(home, ".vvmn")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if _, err := os.Stat(VvmnDir); err != nil {
		if err := os.MkdirAll(VvmnDir, 0777); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	etcDir := filepath.Join(VvmnDir, "etc")
	if _, err := os.Stat(etcDir); err != nil {
		if err := os.Mkdir(etcDir, 0777); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}
	loginFile := filepath.Join(VvmnDir, "etc", "login")
	if _, err := os.Stat(loginFile); err != nil {
		err := ioutil.WriteFile(loginFile, []byte(strings.TrimSpace(`
#!/bin/bash

__vvmn_configure_path()
{
  local vvmn_bin_path="$HOME/.vvmn/vims/current/bin"

  echo "$PATH" | grep -Fqv "$vvmn_bin_path" &&
    PATH="$vvmn_bin_path:$PATH"
}


__vvmn_configure_path

# __END__
		`)), 0666)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
	}

	for _, cmd := range commands {
		if cmd.Name() == args[0] {
			cmd.Flag.Usage = func() { cmd.Usage() }

			cmd.Flag.Parse(args[1:])
			args = cmd.Flag.Args()

			os.Exit(cmd.Run(args))
		}
	}

	fmt.Fprintf(os.Stderr, "vvmn: unknown subcommand %q\nRun ' vvmn help' for usage.\n", args[0])
	os.Exit(2)
}

var usageTemplate = `vvmn is a tool for managing Vim versions

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

func usage() {
	printUsage(os.Stderr)
	os.Exit(2)
}

// help implements the 'help' command.
func help(args []string) {
	if len(args) == 0 {
		printUsage(os.Stdout)
		// not exit 2: succeeded at 'vvmn help'.
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stderr, "usage: vvmn help command\n\nToo many arguments given.\n")
		os.Exit(2) // failed at 'vvmn help'
	}

	arg := args[0]

	for _, cmd := range commands {
		if cmd.Name() == arg {
			tmpl(os.Stdout, helpTemplate, cmd)
			// not exit 2: succeeded at 'vvmn help cmd'.
			return
		}
	}

	fmt.Fprintf(os.Stderr, "Unknown help topic %#q.  Run 'vvmn help'.\n", arg)
	os.Exit(2) // failed at 'vvmn help cmd'
}
