package main

import "fmt"

func (c *Connection) pwd(args []string) {
	if len(args) > 0 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	c.writeout(fmt.Sprintf("257 %q is current directory", c.workdir))
}
