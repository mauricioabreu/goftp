package main

import (
	"fmt"
	"strconv"
	"strings"
)

// Dataport data connection port address
type Dataport struct {
	ip   string
	port int
}

// Address return a valid address with IP+port
func (d *Dataport) Address() string {
	return fmt.Sprintf("%s:%d", d.ip, d.port)
}

func parse(dataport string) (*Dataport, error) {
	fields := strings.Split(dataport, ",")
	ip, sport := fields[0:4], fields[4:]
	port, err := toDecimal(sport)
	if err != nil {
		return nil, err
	}
	return &Dataport{ip: strings.Join(ip, "."), port: port}, nil
}

func toDecimal(rawValues []string) (int, error) {
	fnum, err := strconv.Atoi(rawValues[0])
	if err != nil {
		return 0, err
	}
	snum, err := strconv.Atoi(rawValues[1])
	if err != nil {
		return 0, err
	}
	return (fnum * 256) + snum, nil
}
