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
		case "GET":
			c.get(args)
		case "CLOSE":
			c.writeout("exiting...")
			return
		case "CWD":
			c.cwd(args)
		case "PWD":
			c.pwd()
		default:
			c.writeout("unknown command: %s", command)
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
		c.writeout("bad number of arguments")
		return
	}

	file, err := os.Open(filename)
	if err != nil {
		c.writeout("file not found")
		return
	}

	fileinfo, err := os.Stat(filename)
	if err != nil {
		c.writeout(fmt.Sprintf("could not read file: %s", filename))
		return
	}

	if fileinfo.IsDir() {
		dirs, err := file.Readdirnames(0) // 0 to read all names
		if err != nil {
			c.writeout("error reading dirs inside %s")
			return
		}
		for _, dir := range dirs {
			c.writeout(dir)
		}
	} else {
		c.writeout(filename)
		return
	}
	c.writeout("successful list")
}

func (c *Connection) get(args []string) {
	if len(args) != 1 {
		c.writeout("bad number of arguments")
		return
	}

	filename := args[0]
	file, err := os.Open(filename)
	if err != nil {
		c.writeout("error opening file")
		return
	}

	_, err = io.Copy(c.conn, file)
	if err != nil {
		c.writeout("file unavailable for operation")
		return
	}
	c.writeout("successful get")
}

func (c *Connection) cwd(args []string) {
	if len(args) != 1 {
		c.writeout("bad number of arguments")
		return
	}
	wd := filepath.Join(c.workdir, args[0])
	target := filepath.Join(c.rootdir, wd)

	if _, err := os.Stat(target); err != nil {
		c.writeout("file unavailable for operation")
		return
	}
	c.workdir = wd
	c.writeout("sucessful cwd")
}

func (c *Connection) pwd() {
	c.writeout(filepath.Join(c.rootdir, c.workdir))
}

func (c *Connection) writeout(msg ...interface{}) {
	msg = append(msg, "\r\n")
	fmt.Fprint(c.conn, msg...)
}
