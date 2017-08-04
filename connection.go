package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// Connection : wrapper for connecting and stuff
// TODO rename to FTPCmdConn or sth
type Connection struct {
	conn *net.TCPConn
	buf  *bufio.Reader
}

// NewConnection : constructor for connection
func NewConnection(host, port string) (Connection, error) {
	addr := strings.Join([]string{host, port}, ":")
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return Connection{}, err
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return Connection{}, err
	}

	return Connection{conn: tcpConn, buf: bufio.NewReader(tcpConn)}, nil
}

// ReadLine : read to newline and print the result
func (c *Connection) ReadLine() (line string, err error) {
	if line, err = c.buf.ReadString('\n'); err == nil {
		fmt.Print(line)
	}

	return
}

func (c *Connection) Write(msg string) (int, error) {
	return c.conn.Write([]byte(msg + "\r\n"))
}

// Close : close the connection
func (c *Connection) Close() error {
	return c.conn.Close()
}

// type FTPDataConn struct {
