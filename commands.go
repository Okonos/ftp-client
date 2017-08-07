package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
)

func pwd(cmdConn FTPCmdConn) {
	cmdConn.Exec("PWD")
}

func cd(cmdConn FTPCmdConn, dir string) {
	cmdConn.Exec("CWD " + dir)
}

func ls(cmdConn FTPCmdConn) {
	dataConn, err := cmdConn.NewDataConn() // NewFTPConn(parseHostPort(response))
	if err != nil {
		fmt.Println("Could not initialize data connection: ", err)
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

func quit(cmdConn FTPCmdConn) {
	cmdConn.Exec("QUIT")
}
