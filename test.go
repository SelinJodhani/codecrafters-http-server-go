package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConn(conn net.Conn) {
	for {
		r := bufio.NewReader(conn)
		s, _ := r.ReadString(10) // 10 - \n
		conn.Write([]byte(s))
	}
}

func main() {
	l, err := net.Listen("tcp", ":3333")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleConn(conn)
	}
}
