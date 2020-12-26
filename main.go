package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type Connection struct {
	conn    net.Conn
	rootdir string
	workdir string
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
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatalf("error listening: %s", err)
	}

	for {
		c, err := listener.Accept()
		if err != nil {
			log.Printf("error accepting new connection: %s", err)
			os.Exit(1)
		}
		conn := NewConn(c)
		go conn.handle()
	}
}

func (c *Connection) handle() {
	b := bufio.NewScanner(c.conn)
	var args []string

	for b.Scan() {
		fields := strings.Fields(b.Text())
		if len(fields) == 0 {
			continue
		}
		command := strings.ToUpper(fields[0])
		if len(fields) >= 1 {
			args = fields[1:]
		}

		switch command {
		case "LIST":
			c.list(args)
		case "RETR":
			c.retr(args)
		case "QUIT":
			c.writeout("221 Service closing control connection.")
			return
		case "CWD":
			c.cwd(args)
		case "PWD":
			c.pwd()
		default:
			c.writeout("502 command not implemented.")
			continue
		}
	}

	if b.Err() != nil {
		log.Println(b.Err())
	}
}

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

	fileinfo, err := os.Stat(filename)
	if err != nil {
		log.Println(fmt.Sprintf("could not read file: %s", filename))
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}

	if fileinfo.IsDir() {
		dirs, err := file.Readdirnames(0) // 0 to read all names
		if err != nil {
			c.writeout("450 Requested file action not taken.")
			return
		}
		for _, dir := range dirs {
			c.writeout(dir)
		}
	} else {
		c.writeout(filename)
		return
	}
	c.writeout("200 successful command.")
}

func (c *Connection) retr(args []string) {
	if len(args) != 1 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}

	filename := args[0]
	file, err := os.Open(filename)
	if err != nil {
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}

	_, err = io.Copy(c.conn, file)
	if err != nil {
		c.writeout("450 Requested file action not taken.")
		return
	}
	c.writeout("200 successful command.")
}

func (c *Connection) cwd(args []string) {
	if len(args) != 1 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	wd := filepath.Join(c.workdir, args[0])
	target := filepath.Join(c.rootdir, wd)

	if _, err := os.Stat(target); err != nil {
		c.writeout("550 Requested action not taken. File unavailable.")
		return
	}
	c.workdir = wd
	c.writeout("200 successful command.")
}

func (c *Connection) pwd() {
	c.writeout(filepath.Join(c.rootdir, c.workdir))
}

func (c *Connection) writeout(msg ...interface{}) {
	msg = append(msg, "\r\n")
	fmt.Fprint(c.conn, msg...)
}
