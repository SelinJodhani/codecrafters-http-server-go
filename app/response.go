package main

import (
	"fmt"
	"net"
)

func SendResponse(conn net.Conn, statusCode int, data string) {
	conn.Write(
		[]byte(
			fmt.Sprintf(
				"HTTP/1.1 %d OK\r\nContent-Type: text/plain\r\nContent-Length: %d\r\n\n%s\r\n",
				statusCode,
				len(data),
				data,
			),
		),
	)
}
