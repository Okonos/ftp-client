package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
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

func get(cmdConn FTPCmdConn, filename string) {
	dataConn, err := cmdConn.NewDataConn()
	if err != nil {
		fmt.Println("Could not initialize data connection: ", err)
		return
	}
	defer dataConn.Close()

	resp, err := cmdConn.Exec("RETR " + filename)
	if err != nil {
		fmt.Println("Error sending command: ", err)
		return
	}
	// check for negative reply
	if resp[0] == '4' || resp[0] == '5' {
		return
	}

	f, err := os.Create(filename)
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return
	}
	defer f.Close()

	buf := make([]byte, 1024)
	for {
		if _, err := dataConn.Read(buf); err != nil {
			if err != io.EOF {
				fmt.Println("Error reading response: ", err)
				return
			}
			break
		}

		if _, err := f.Write(buf); err != nil {
			fmt.Println("Error writing to file: ", err)
			return
		}
	}

	cmdConn.ReadLine()
}

func quit(cmdConn FTPCmdConn) {
	cmdConn.Exec("QUIT")
}
