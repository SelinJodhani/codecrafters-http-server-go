package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

type HTTPServer struct {
	listener  net.Listener
	directory string
}

func (s *HTTPServer) Close() {
	s.listener.Close()
}

func (s *HTTPServer) Serve() {
	defer s.Close()
	fmt.Println("Server listening on port", s.listener.Addr().(*net.TCPAddr).Port)

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go s.HandleConnection(conn)
	}
}

func (s *HTTPServer) readInput(c net.Conn) (string, error) {
	reader := bufio.NewReader(c)
	var requestBuffer bytes.Buffer

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		requestBuffer.WriteString(line)
		if line == "\r\n" {
			break // Empty line indicates end of headers
		}
	}

	return requestBuffer.String(), nil
}

func (s *HTTPServer) HandleConnection(c net.Conn) {
	defer c.Close()

	requestString, err := s.readInput(c)
	if err != nil {
		fmt.Println("Error reading request:", err)
		return
	}

	request, err := ParseHTTPRequest(requestString)
	if err != nil {
		fmt.Println("Error parsing request:", err)
		return
	}

	// Check if file exists

	switch {
	case request.Path == "/":
		SendResponse(c, 200, "")
	case strings.HasPrefix(request.Path, "/echo"):
		response := strings.TrimPrefix(request.Path, "/echo/")
		SendResponse(c, 200, response)
	case strings.HasPrefix(request.Path, "/user-agent"):
		userAgent := request.Headers["User-Agent"]
		SendResponse(c, 200, userAgent)
	case strings.HasPrefix(request.Path, "/files"):
		fileName := strings.TrimPrefix(request.Path, "/files/")
		filePath := filepath.Join(s.directory, fileName)

		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			SendResponse(c, 404, "")
			return
		}

		fileContents, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			SendResponse(c, 500, "")
			return
		}

		c.Write(
			[]byte(
				fmt.Sprintf(
					"HTTP/1.1 200 OK\r\nContent-Type: application/octet-stream\r\nContent-Length: %d\r\n\n%s\r\n",
					len(string(fileContents)),
					string(fileContents),
				),
			),
		)
	default:
		SendResponse(c, 404, "")
	}
}

func NewHTTPServer(port string, directory string) (*HTTPServer, error) {
	listener, err := net.Listen("tcp", port)

	if err != nil {
		return nil, err
	}

	return &HTTPServer{
		listener:  listener,
		directory: directory,
	}, nil
}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	dicrectory := flag.String("directory", ".", "Specify the directory")
	flag.Parse()

	server, err := NewHTTPServer(":4221", *dicrectory)

	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}

	defer server.Close()

	server.Serve()
}
