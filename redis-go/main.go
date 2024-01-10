package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error with connection: %s\n", err)
			continue
		}

		go handleConnection(conn)

	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading msg: %s\n", err)
			return
		}
		resp := fmt.Sprintf("Got message: %s\n", msg)
		resp = strings.TrimSpace(resp)
		args := parseArgs(msg)
		handleArgs(args, conn)

	}

}

func parseArgs(args string) []string {
	// TODO make this smarter to handle quoted strings
	split := strings.Split(args, " ")
	for i, arg := range split {
		split[i] = strings.TrimSpace(arg)
	}
	return split
}

func ping(str string, conn net.Conn) {
	var err error
	if len(str) > 0 {
		_, err = conn.Write([]byte(fmt.Sprintf("PONG %s\n", str)))
	} else {
		_, err = conn.Write([]byte("PONG\n"))
	}
	if err != nil {
		fmt.Printf("Error writing: %s", err)
	}
}

func echo(args []string, conn net.Conn) {
	writeString := ""
	if len(args) > 0 {
		writeString = strings.Join(args, " ") + "\n"
	} else {
		writeString = "unsupported, ECHO requires [1+] arguments"
	}
	_, err := conn.Write([]byte(writeString))
	if err != nil {
		fmt.Printf("Error writing: %s", err)
	}
}

func handleArgs(args []string, conn net.Conn) {
	switch strings.ToLower(args[0]) {
	case "ping":
		if len(args) > 1 {
			ping(args[1], conn)
		} else {
			ping("", conn)
		}
	case "echo":
		if len(args) > 1 {
			echo(args[1:], conn)
		} else {
			echo([]string{}, conn)
		}

	case "llen":
		fmt.Printf("LLEN not implemented yet\n")
	default:
		fmt.Printf("unsupported command: %s\n", args[0])

	}
}
