package main

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/stefanhts/redis-mini/data"
)

var store *data.Store

func main() {
	store = data.NewStore()
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

func errorMsg(conn net.Conn, cmd string) {
	_, err := conn.Write([]byte("Error running command: " + cmd))
	if err != nil {
		fmt.Printf("Error writing: %s", err)
	}
}

func push(conn net.Conn, key string, vals ...string) {
	pushed := store.LPush(key, vals...)
	msg := fmt.Sprintf("Pushed %d values to key %s \n", pushed, key)
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Error writing: %s\n", err)
	}
}

func pop(conn net.Conn, key string, num int64) {
	els, err := store.LPop(key, num)
	msg := "Popped: "

	if err != nil {
		msg = fmt.Sprintf("Error popping from list: %s\n", err)
	}

	for _, el := range els {
		msg += el + " "
	}

	if len(els) == 0 {
		msg += "nothing"
	}

	msg += "from key " + key + "\n"

	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Error writing: %s\n", err)
	}
}

func llen(conn net.Conn, key string) {
	length := store.LLen(key)
	msg := fmt.Sprintf("LLEN: %d, for key: %s\n", length, key)
	_, err := conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Error writing: %s\n", err)
	}
}

func lpos(conn net.Conn, key string, val string) {
	pos, err := store.LPos(key, val)
	msg := fmt.Sprintf("Position for key: %s, val: %s is: %d\n", key, val, pos)
	if err != nil {
		msg = "error getting LPOS\n"
	}
	_, err = conn.Write([]byte(msg))
	if err != nil {
		fmt.Printf("Error writing: %s\n", err)
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
	case "lpush":
		if len(args) < 3 {
			errorMsg(conn, args[0])
		} else {
			push(conn, args[1], args[2:]...)
		}
	case "lpop":
		if len(args) <= 1 {
			errorMsg(conn, args[0])
		} else {
			num, _ := strconv.ParseInt(args[2], 10, 64)
			pop(conn, args[1], num)
		}
	case "llen":
		if len(args) <= 1 {
			errorMsg(conn, args[0])
		} else {
			llen(conn, args[1])
		}

	case "lpos":
		if len(args) != 3 {
			errorMsg(conn, args[0])
		} else {
			lpos(conn, args[1], args[2])
		}

	default:
		fmt.Printf("unsupported command: %s\n", args[0])

	}
}
