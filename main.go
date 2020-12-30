package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"path/filepath"
	"strings"
)

var port int

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

func main() {
	flag.IntVar(&port, "port", 8080, "port to run server")
	flag.Parse()

	listener, err := net.Listen("tcp4", fmt.Sprintf(":%d", port))
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

	c.writeout("220 Welcome")
	for b.Scan() {
		fields := strings.Fields(b.Text())
		if len(fields) == 0 {
			continue
		}
		command := strings.ToUpper(fields[0])
		if len(fields) >= 1 {
			args = fields[1:]
		}

		log.Printf("[command] %s %v\n", command, args)
		switch command {
		case "LIST":
			c.list(args)
		case "RETR":
			c.retr(args)
		case "STOR":
			c.stor(args)
		case "QUIT":
			c.writeout("221 Service closing control connection.")
			return
		case "CWD":
			c.cwd(args)
		case "PWD":
			c.pwd(args)
		case "PORT":
			c.port(args)
		case "PASV":
			c.pasv(args)
		case "NOOP":
			c.noop(args)
		case "TYPE":
			c.setType(args)
		case "STRU":
			c.stru(args)
		case "MODE":
			c.mode(args)
		case "USER":
			c.user(args)
		default:
			c.writeout("502 command not implemented.")
			continue
		}
		if command != "PASV" && c.listener != nil {
			c.cleanupListener()
		}
		c.prevCommand = command
	}

	if b.Err() != nil {
		log.Println(b.Err())
	}
}

func (c *Connection) port(args []string) {
	if len(args) != 1 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}

	address := args[0]
	dataport, err := parse(address)
	if err != nil {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	c.dataport = dataport
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

func (c *Connection) stor(args []string) {
	if len(args) != 1 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}

	filename := filepath.Join(c.rootdir, c.workdir, args[0])
	file, err := os.Create(filename)
	if err != nil {
		log.Println(err)
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

	_, err = io.Copy(file, dc)
	if err != nil {
		log.Println(err)
		c.writeout("450 Requested file action not taken.")
		return
	}
	c.writeout("226 Closing data connection. Requested file action successful.")
}

func (c *Connection) noop(args []string) {
	if len(args) > 0 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	c.writeout("200 successful command.")
}

func (c *Connection) setType(args []string) {
	flag := strings.ToUpper(strings.Join(args, ""))
	switch flag {
	case "A", "A N":
		c.binary = false
	case "I", "L 8":
		c.binary = true
	default:
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}
	c.writeout("200 successful command.")
}

func (c *Connection) stru(args []string) {
	value := strings.Join(args, "")
	if strings.EqualFold(value, "F") {
		c.writeout("200 successful command.")
		return
	}
	c.writeout("504 Command not implemented for that parameter.")
}

func (c *Connection) mode(args []string) {
	value := strings.Join(args, "")
	if strings.EqualFold(value, "S") {
		c.writeout("200 successful command.")
		return
	}
	c.writeout("504 Command not implemented for that parameter.")
}

func (c *Connection) user(args []string) {
	c.writeout("200 successful command.")
}

func (c *Connection) lineterminator() string {
	if c.binary {
		return "\n"
	}
	return "\r\n"
}

func (c *Connection) writeout(msg ...interface{}) {
	msg = append(msg, c.lineterminator())
	fmt.Fprint(c.conn, msg...)
}
