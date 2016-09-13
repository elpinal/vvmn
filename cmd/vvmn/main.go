package main

import (
	"fmt"
	"io"
	"os"

	"github.com/susp/vvmn/cli"
)

type colorWriter struct {
	w     io.Writer
	color int
}

func (c *colorWriter) Write(b []byte) (int, error) {
	var n int
	if _, err := io.WriteString(c.w, fmt.Sprintf("\033[%vm", c.color)); err != nil {
		return n, err
	}
	n0, err := c.w.Write(b)
	n += n0
	if err != nil {
		return n, err
	}
	if _, err := io.WriteString(c.w, "\033[0m"); err != nil {
		return n, err
	}
	return len(b), nil
}

func main() {
	yw := &colorWriter{w: os.Stdout, color: 33}
	rw := &colorWriter{w: os.Stdout, color: 31}
	c := cli.CLI{OutStream: yw, ErrStream: rw}
	r := c.Run(os.Args[1:])
	os.Exit(r)
}
