package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
)

func ls(cmdConn FTPCmdConn) {
	buf, _ := cmdConn.Exec("PASV")

	dataConn, err := NewFTPConn(parseHostPort(buf))
	if err != nil {
		fmt.Println("Could not initialize connection: ", err)
		return
	}
	defer dataConn.Close()

	cmdConn.Exec("LIST")
	dataReader := bufio.NewReader(dataConn)
	data, err := ioutil.ReadAll(dataReader)
	if err != nil {
		fmt.Println("Error reading response: ", err)
		return
	}
	fmt.Print(string(data))
	cmdConn.ReadLine()
}
