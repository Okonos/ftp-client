package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func parseArgs() (addr string) {
	args := os.Args[1:]
	if len(args) > 0 && args[0] != "" {
		host := args[0]
		port := "21"
		if len(args) > 1 && args[1] != "" {
			port = args[1]
		}
		addr = fmt.Sprintf("%s:%s", host, port)
	}

	return
}

func readInput(prompt string) (text string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err = reader.ReadString('\n')
	if err == nil {
		text = text[:len(text)-1] // cut the '\n'
	}
	return
}

func promptLoop() {
	for {
		cmd, err := readInput("ftp> ")
		if err != nil {
			if err == io.EOF {
				// TODO disconnect
				return
			}
			log.Fatal(err)
		}

		switch cmd {
		default:
			fmt.Println("Not implemented")
		}
	}
}

func main() {
	addr := parseArgs()
	if addr == "" {
		log.Fatal("Host address must be specified")
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	cmdConn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer cmdConn.Close()

	fmt.Println("Connected")

	bufReader := bufio.NewReader(cmdConn)
	buf, _ := bufReader.ReadString('\n')
	fmt.Println(buf)

	if name, err := readInput("Name: "); err == nil {
		cmdConn.Write([]byte("USER " + name + "\r\n"))
		buf, _ = bufReader.ReadString('\n')
		fmt.Println(buf)

		if pass, err := readInput("Password: "); err == nil || err == io.EOF {
			cmdConn.Write([]byte("PASS " + pass + "\r\n"))
			buf, _ = bufReader.ReadString('\n')
			fmt.Println(buf)
		} else {
			log.Fatal(err)
		}
	} else {
		if err == io.EOF {
			fmt.Println("Login failed.")
		} else {
			log.Fatal(err)
		}
	}

	promptLoop()
}
