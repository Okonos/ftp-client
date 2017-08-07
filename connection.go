package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// FTPConn : command connection wrapper struct
type FTPConn struct {
	conn *net.TCPConn
	buf  *bufio.Reader
}

// NewFTPConn : constructor for connection
func NewFTPConn(host, port string) (*FTPConn, error) {
	addr := strings.Join([]string{host, port}, ":")
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return &FTPConn{}, err
	}
	tcpConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return &FTPConn{}, err
	}

	return &FTPConn{conn: tcpConn, buf: bufio.NewReader(tcpConn)}, nil
}

func (c *FTPConn) Read(b []byte) (n int, err error) {
	return c.conn.Read(b)
}

// ReadLine : read to newline and print the result
func (c *FTPConn) ReadLine() (line string, err error) {
	if line, err = c.buf.ReadString('\n'); err == nil {
		fmt.Print(line)
	}

	return
}

func (c *FTPConn) Write(msg string) (int, error) {
	return c.conn.Write([]byte(msg + "\r\n"))
}

// Exec : Write command and read the response
func (c *FTPConn) Exec(cmd string) (line string, err error) {
	c.Write(cmd)
	return c.ReadLine()
}

// Close : close the connection
func (c *FTPConn) Close() error {
	return c.conn.Close()
}

// type FTPDataConn struct {
