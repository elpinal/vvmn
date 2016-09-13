package main

import (
	"os"

	"github.com/susp/vvmn/cli"
)

func main() {
	c := cli.CLI{OutStream: os.Stdout, ErrStream: os.Stderr}
	r := c.Run(os.Args[1:])
	os.Exit(r)
}
