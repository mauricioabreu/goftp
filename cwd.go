package main

import (
	"log"
	"os"
	"path/filepath"
)

func (c *Connection) cwd(args []string) {
	if len(args) != 1 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	wd := filepath.Join(c.workdir, args[0])
	target := filepath.Join(c.rootdir, wd)

	if _, err := os.Stat(target); err != nil {
		log.Println(err)
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}
	c.workdir = wd
	c.writeout("200 successful command.")
}
