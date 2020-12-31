package main

import (
	"os"
	"path/filepath"
)

func (c *Connection) dele(args []string) {
	if len(args) != 1 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	filename := filepath.Join(c.curDir(), filepath.Clean(args[0]))
	err := os.Remove(filename)
	if err != nil {
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}
	c.writeout("250 Requested file action okay, completed.")
}
