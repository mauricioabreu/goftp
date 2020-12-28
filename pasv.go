package main

import (
	"fmt"
	"net"
	"strconv"
)

func (c *Connection) pasv(args []string) {
	if len(args) > 0 {
		c.writeout("501 Syntax error in parameters or arguments.")
		return
	}

	var err error
	c.listener, err = net.Listen("tcp4", "") // empty address for automatically choosen port
	if err != nil {
		c.cleanupListener()
		c.writeout("415 Requested action aborted. Local error in processing.")
		return
	}
	_, port, err := net.SplitHostPort(c.listener.Addr().String())
	if err != nil {
		c.cleanupListener()
		c.writeout("415 Requested action aborted. Local error in processing.")
		return
	}
	ip, _, err := net.SplitHostPort(c.conn.LocalAddr().String())
	if err != nil {
		c.cleanupListener()
		c.writeout("415 Requested action aborted. Local error in processing.")
		return
	}
	addr, err := toFTPAddress(ip, port)
	if err != nil {
		c.cleanupListener()
		c.writeout("415 Requested action aborted. Local error in processing.")
		return
	}
	c.writeout(fmt.Sprintf("227 =%s", addr))
}

func (c *Connection) cleanupListener() {
	c.listener.Close()
	c.listener = nil
}

// Convert a host/port address to FTP address style
// Used to answer a PASV command. Example: 227 =h1,h2,h3,h4,p1,p2
func toFTPAddress(host, sport string) (string, error) {
	ipAddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return "", err
	}
	port, err = strconv.Atoi(sport)
	if err != nil {
		return "", err
	}
	ipRepr := ipAddr.IP.To4()
	return fmt.Sprintf("%d,%d,%d,%d,%d,%d", ipRepr[0], ipRepr[1], ipRepr[2], ipRepr[3], port/256, port%256), nil
}
