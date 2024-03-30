package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	r := bufio.NewReader(conn)

	str, err := r.ReadString('\n')
	if err != nil {
		panic(err)
	}

	line := strings.Split(str, " ")
	path := line[1]

	if path == "/" {
		conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else if strings.HasPrefix(path, "/echo") {
		e := strings.Split(path, "/")
		randomStr := e[2:]
		str := strings.Join(randomStr, "/")
		conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\n%s\r\n", len(str), str)))
	} else if strings.HasPrefix(path, "/user-agent") {
		for {
			str, err := r.ReadString('\n')
			if err != nil {
				fmt.Println("ERR: ", err)
				break
			}
			if strings.HasPrefix(str, "User-Agent") {
				agent := strings.Split(str, ":")
				text := strings.TrimSpace(agent[1])
				conn.Write([]byte(fmt.Sprintf("HTTP/1.1 200 OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\n%s\r\n", len(text), text)))
				break
			}
		}
	} else {
		conn.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
