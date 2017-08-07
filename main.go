package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func readInput(prompt string) (text string, err error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	text, err = reader.ReadString('\n')
	if err == nil {
		text = text[:len(text)-1] // cut the '\n'
	}
	return
}

func parseArgs() (host, port string) {
	args := os.Args[1:]
	if len(args) > 0 && args[0] != "" {
		host = args[0]
		port = "21"
		if len(args) > 1 && args[1] != "" {
			port = args[1]
		}
	}

	return
}

func parseHostPort(addr string) (host, port string) {
	start, end := strings.Index(addr, "(")+1, strings.Index(addr, ")")
	addrBytes := strings.Split(addr[start:end], ",")
	host = strings.Join(addrBytes[:4], ".")
	var portVal int
	if upperByte, err := strconv.Atoi(addrBytes[4]); err == nil {
		portVal += upperByte * 256
	}
	if lowerByte, err := strconv.Atoi(addrBytes[5]); err == nil {
		portVal += lowerByte
	}
	port = strconv.Itoa(portVal)
	return
}

func promptLoop(cmdConn FTPCmdConn) {
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
		case "ls":
			ls(cmdConn)
		default:
			fmt.Println("Invalid command")
		}
	}
}

func main() {
	host, port := parseArgs()
	if host == "" {
		log.Fatal("Host address must be specified")
	}
	cmdConn, err := NewFTPConn(host, port)
	if err != nil {
		log.Fatal(err)
	}
	defer cmdConn.Close()

	fmt.Println("Connected")

	cmdConn.ReadLine()

	if name, err := readInput("Name: "); err == nil {
		cmdConn.Exec("USER " + name)

		if pass, err := readInput("Password: "); err == nil || err == io.EOF {
			cmdConn.Exec("PASS " + pass)
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

	promptLoop(cmdConn)
}
