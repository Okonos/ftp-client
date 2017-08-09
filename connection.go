package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// FTPConn : command connection wrapper struct
type FTPConn struct {
	conn *net.TCPConn
	buf  *bufio.Reader
}

// FTPCmdConn : FTP command connection interface
type FTPCmdConn interface {
	Read([]byte) (int, error)
	ReadLine() (string, error)
	Write([]byte) (int, error)
	WriteCmd(string) (int, error)
	Exec(string) (string, error)
	InitDataConn() (*FTPConn, error)
	Close() error
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

// InitDataConn : establish data connection
func (c *FTPConn) InitDataConn() (*FTPConn, error) {
	response, err := c.Exec("PASV")
	if err != nil {
		return &FTPConn{}, err
	}

	start, end := strings.Index(response, "(")+1, strings.Index(response, ")")
	addr := strings.Split(response[start:end], ",")
	host := strings.Join(addr[:4], ".")
	var portVal int
	if upperByte, err := strconv.Atoi(addr[4]); err == nil {
		portVal += upperByte * 256
	}
	if lowerByte, err := strconv.Atoi(addr[5]); err == nil {
		portVal += lowerByte
	}
	port := strconv.Itoa(portVal)
	return NewFTPConn(host, port)
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

func (c *FTPConn) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

// WriteCmd : write string with crlf appended
func (c *FTPConn) WriteCmd(msg string) (int, error) {
	return c.conn.Write([]byte(msg + "\r\n"))
}

// Exec : Write command and read the response
func (c *FTPConn) Exec(cmd string) (line string, err error) {
	c.WriteCmd(cmd)
	return c.ReadLine()
}

// Close : close the connection
func (c *FTPConn) Close() error {
	return c.conn.Close()
}
