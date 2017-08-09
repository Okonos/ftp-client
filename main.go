package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
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

func readPassword() (string, error) {
	fmt.Print("Password: ")
	bytePasswd, err := terminal.ReadPassword(int(syscall.Stdin))
	return strings.TrimSpace(string(bytePasswd)), err
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

func promptLoop(cmdConn FTPCmdConn) {
	for {
		input, err := readInput("ftp> ")
		if err != nil {
			if err == io.EOF {
				quit(cmdConn)
				return
			}
			log.Fatal(err)
		}

		var cmd, arg string
		if len(input) > 0 {
			splitInput := strings.Split(input, " ")
			cmd = splitInput[0]
			if len(splitInput) > 1 {
				arg = splitInput[1]
			}
		}

		switch strings.ToLower(cmd) {
		case "pwd":
			pwd(cmdConn)
		case "cd":
			cd(cmdConn, arg)
		case "ls":
			ls(cmdConn)
		case "get":
			get(cmdConn, arg)
		case "exit", "quit":
			quit(cmdConn)
			return
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
		if pass, err := readPassword(); err == nil || err == io.EOF {
			fmt.Println()
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

	// binary mode
	cmdConn.Exec("TYPE I")

	promptLoop(cmdConn)
}
