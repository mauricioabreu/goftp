package main

import (
	"errors"
	"io"
	"net"
)

// ErrBadSequence bad sequence of commands
var ErrBadSequence = errors.New("bad sequence of commands")

func (c *Connection) dataconnection() (io.ReadWriteCloser, error) {
	if c.prevCommand != "PORT" {
		return nil, ErrBadSequence
	}
	conn, err := net.Dial("tcp", c.dataport.Address())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
