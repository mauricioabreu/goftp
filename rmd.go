package main

import (
	"os"
	"path/filepath"
)

func (c *Connection) rmd(args []string) {
	if len(args) != 1 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	dir := filepath.Join(c.curDir(), filepath.Clean(args[0]))
	if fileinfo, err := os.Stat(dir); err == nil && !fileinfo.IsDir() {
		c.writeout("450 Requested file action not taken.")
		return
	}
	err := os.RemoveAll(dir)
	if err != nil {
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}
	c.writeout("250 Requested file action okay, completed.")
}
