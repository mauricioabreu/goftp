package main

import (
	"os"
	"path/filepath"
)

func (c *Connection) mkd(args []string) {
	if len(args) != 1 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	dir := filepath.Join(c.curDir(), filepath.Clean(args[0]))
	err := os.Mkdir(dir, 0755)
	if err != nil {
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}
	c.writeout("257 %q directory created.", filepath.Base(dir))
}
