package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
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

func parseHostPort(addr string) (*net.TCPAddr, error) {
	start, end := strings.Index(addr, "(")+1, strings.Index(addr, ")")
	addrBytes := strings.Split(addr[start:end], ",")
	host := strings.Join(addrBytes[:4], ".")
	var port int
	if upperByte, err := strconv.Atoi(addrBytes[4]); err == nil {
		port += upperByte * 256
	}
	if lowerByte, err := strconv.Atoi(addrBytes[5]); err == nil {
		port += lowerByte
	}
	fullAddr := strings.Join([]string{host, strconv.Itoa(port)}, ":")
	return net.ResolveTCPAddr("tcp", fullAddr)
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
		case "ls":
			return
		default:
			fmt.Println("Invalid command")
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
	fmt.Print(buf)

	if name, err := readInput("Name: "); err == nil {
		cmdConn.Write([]byte("USER " + name + "\r\n"))
		buf, _ = bufReader.ReadString('\n')
		fmt.Print(buf)

		if pass, err := readInput("Password: "); err == nil || err == io.EOF {
			cmdConn.Write([]byte("PASS " + pass + "\r\n"))
			buf, _ = bufReader.ReadString('\n')
			fmt.Print(buf)
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

	cmdConn.Write([]byte("PASV\r\n"))
	buf, _ = bufReader.ReadString('\n')
	fmt.Print(buf)

	dataTCPAddr, err := parseHostPort(buf)
	if err != nil {
		log.Fatal(err)
	}
	dataConn, err := net.DialTCP("tcp", nil, dataTCPAddr)
	if err != nil {
		log.Fatal(err)
	}
	cmdConn.Write([]byte("LIST\r\n"))
	buf, _ = bufReader.ReadString('\n')
	fmt.Print(buf)
	dataReader := bufio.NewReader(dataConn)
	data, _ := ioutil.ReadAll(dataReader)
	fmt.Print(string(data))
	dataConn.Close() // defer
	buf, _ = bufReader.ReadString('\n')
	fmt.Print(buf)

	promptLoop()
}
