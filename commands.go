package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

func pwd(cmdConn FTPCmdConn) {
	cmdConn.Exec("PWD")
}

func cd(cmdConn FTPCmdConn, dir string) {
	cmdConn.Exec("CWD " + dir)
}

func ls(cmdConn FTPCmdConn) {
	dataConn, err := cmdConn.InitDataConn()
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
	dataConn, err := cmdConn.InitDataConn()
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

	buf := make([]byte, 8192)
	var bytesRcvd int
	start := time.Now()

	for {
		n, err := dataConn.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading response: ", err)
			return
		}

		if _, err := f.Write(buf); err != nil {
			fmt.Println("Error writing to file: ", err)
			return
		}

		bytesRcvd += n
	}

	secs := float64(time.Now().Sub(start)) / float64(time.Second)

	cmdConn.ReadLine()
	fmt.Printf("%d bytes received in %.2f secs (%.4f MB/s)\n",
		bytesRcvd, secs, float64(bytesRcvd/1024/1024)/secs)
}

func quit(cmdConn FTPCmdConn) {
	cmdConn.Exec("QUIT")
}
