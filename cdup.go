package main

import (
	"log"
	"os"
	"path/filepath"
)

func (c *Connection) cdup(args []string) {
	if len(args) != 0 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	wd := filepath.Join(c.workdir, "..")

	if _, err := os.Stat(filepath.Join(c.rootdir, wd)); err != nil {
		log.Println(err)
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}
	c.workdir = wd
	c.writeout("200 successful command.")
}
