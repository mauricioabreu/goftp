package main

import (
	"log"
	"net"
	"path/filepath"
)

// Connection represent a connection to accept commands
type Connection struct {
	conn        net.Conn
	dataport    *Dataport
	rootdir     string
	workdir     string
	binary      bool
	prevCommand string
	listener    net.Listener
}

// NewConn prepare a connection to be used
func NewConn(c net.Conn) *Connection {
	rd, err := filepath.Abs(".")
	if err != nil {
		log.Fatal(err)
	}
	return &Connection{
		conn:    c,
		rootdir: rd,
		workdir: "/",
		binary:  false,
	}
}
