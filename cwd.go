package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func (c *Connection) cwd(args []string) {
	if len(args) != 1 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	wd := c.buildWorkDir(args[0])
	log.Println(wd)

	if _, err := os.Stat(filepath.Join(c.rootdir, wd)); err != nil {
		log.Println(err)
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}
	c.workdir = wd
	c.writeout("200 successful command.")
}

func (c *Connection) buildWorkDir(path string) string {
	cpath := filepath.Clean(path)
	if cpath[0:1] == "/" {
		if strings.HasPrefix(cpath, c.rootdir) {
			return strings.TrimPrefix(cpath, c.rootdir)
		}
	}
	return filepath.Join(c.workdir, cpath)
}
