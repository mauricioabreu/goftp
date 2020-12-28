package main

import (
	"errors"
	"io"
	"net"
)

// ErrBadSequence bad sequence of commands
var ErrBadSequence = errors.New("bad sequence of commands")

func (c *Connection) dataconnection() (io.ReadWriteCloser, error) {
	if c.prevCommand != "PORT" && c.prevCommand != "PASV" {
		return nil, ErrBadSequence
	}
	var conn net.Conn
	var err error
	if c.prevCommand == "PORT" {
		conn, err = net.Dial("tcp", c.dataport.Address())
		if err != nil {
			return nil, err
		}
	} else {
		conn, err = c.listener.Accept()
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}
