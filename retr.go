package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"path/filepath"
)

func (c *Connection) retr(args []string) {
	if len(args) != 1 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}

	filename := filepath.Join(c.curDir(), filepath.Clean(args[0]))
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

	if c.binary {
		_, err = io.Copy(dc, file)
		if err != nil {
			log.Println(err)
			c.writeout("450 Requested file action not taken.")
			return
		}
		c.writeout("226 Closing data connection. Requested file action successful.")
	} else {
		r, w := bufio.NewReader(file), bufio.NewWriter(dc)
		for {
			line, isPrefix, err := r.ReadLine()
			if err == io.EOF {
				if err := w.Flush(); err != nil {
					log.Println(err)
				}
				c.writeout("226 Closing data connection. Requested file action successful.")
				return
			}
			if err != nil {
				log.Println(err)
				c.writeout("426 Connection closed; transfer aborted.")
				return
			}
			if _, err = w.Write(line); err != nil {
				log.Println(err)
				c.writeout("426 Connection closed; transfer aborted.")
				return
			}
			if !isPrefix {
				if _, err := w.Write([]byte(c.lineterminator())); err != nil {
					log.Println(err)
					c.writeout("426 Connection closed; transfer aborted.")
					return
				}
			}
		}
	}
}
