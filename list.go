package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func (c *Connection) list(args []string) {
	var filename string
	switch lenargs := len(args); lenargs {
	case 0:
		filename = filepath.Join(c.rootdir, c.workdir)
	case 1:
		filename = filepath.Join(c.rootdir, c.workdir, args[0])
	default:
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}

	c.writeout("150 File status okay; about to open data connection.")
	dc, err := c.dataconnection()
	if err == ErrBadSequence {
		c.writeout("503 Bad sequence of commands.")
		return
	}
	if err != nil {
		log.Println(err)
		c.writeout("425 Can't open data connection.")
		return
	}
	defer dc.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		log.Println(fmt.Sprintf("Could not read file: %s", filename))
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}

	if fileinfo.IsDir() {
		files, err := file.Readdirnames(0) // 0 to read all names
		if err != nil {
			log.Println(err)
			c.writeout("450 Requested file action not taken.")
			return
		}
		for _, file := range files {
			if _, err := fmt.Fprint(dc, file, c.lineterminator()); err != nil {
				log.Println(err)
				c.writeout("426 Connection closed; transfer aborted.")
				return
			}
		}
	} else {
		if _, err := fmt.Fprint(dc, filename, c.lineterminator()); err != nil {
			log.Println(err)
			c.writeout("426 Connection closed; transfer aborted.")
			return
		}
	}
	c.writeout("226 Closing data connection. Requested file action successful.")
}
