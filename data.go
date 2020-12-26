package main

import (
	"io"
	"net"
)

func (c *Connection) dataconnection() (io.ReadWriteCloser, error) {
	conn, err := net.Dial("tcp", c.dataport.Address())
	if err != nil {
		return nil, err
	}
	return conn, nil
}
