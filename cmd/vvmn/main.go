package main

import (
	"os"

	"github.com/elpinal/color"
	"github.com/susp/vvmn/cli"
)

func main() {
	yw := color.New(os.Stdout, color.Yellow)
	rw := color.New(os.Stdout, color.Red)
	c := cli.CLI{OutStream: yw, ErrStream: rw}
	r := c.Run(os.Args[1:])
	os.Exit(r)
}
