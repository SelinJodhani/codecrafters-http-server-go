package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
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

	var (
		requestBuffer bytes.Buffer
		bodyLength    int
		hasBody       bool
	)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}

		requestBuffer.WriteString(line)
		if line == "\r\n" {
			break // Empty line indicates end of headers
		}

		if strings.HasPrefix(line, "Content-Length:") {
			parts := strings.Fields(line)
			if len(parts) == 2 {
				bodyLength, _ = strconv.Atoi(parts[1])
				hasBody = true
			}
		}
	}

	// Read request body if present
	if hasBody && bodyLength > 0 {
		bodyBytes := make([]byte, bodyLength)
		_, err := io.ReadFull(reader, bodyBytes)
		if err != nil {
			return "", err
		}
		requestBuffer.Write(bodyBytes)
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

	response := NewResponse()

	if encoding, ok := request.Headers["Accept-Encoding"]; ok &&
		strings.Contains(encoding, "gzip") {
		response.SetHeader("Content-Encoding", "gzip")
	}

	switch {

	case request.Path == "/":
		response.SetStatusCode(200).Send(c)

	case strings.HasPrefix(request.Path, "/echo"):
		message := strings.TrimPrefix(request.Path, "/echo/")
		response.SetStatusCode(200).SetBody(message).Send(c)

	case strings.HasPrefix(request.Path, "/user-agent"):
		userAgent := request.Headers["User-Agent"]
		response.SetStatusCode(200).SetBody(userAgent).Send(c)

	case request.Method == "POST" && strings.HasPrefix(request.Path, "/files"):
		fileName := strings.TrimPrefix(request.Path, "/files/")
		filePath := filepath.Join(s.directory, fileName)

		err := os.WriteFile(filePath, []byte(request.Body), 0644)
		if err != nil {
			fmt.Println("Error saving file:", err)
			response.SetStatusCode(500).Send(c)
			return
		}

		response.SetStatusCode(201).Send(c)

	case request.Method == "GET" && strings.HasPrefix(request.Path, "/files"):
		fileName := strings.TrimPrefix(request.Path, "/files/")
		filePath := filepath.Join(s.directory, fileName)

		_, err = os.Stat(filePath)
		if os.IsNotExist(err) {
			response.SetStatusCode(404).Send(c)
			return
		}

		fileContents, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			response.SetStatusCode(500).Send(c)
			return
		}

		response.SetStatusCode(200).
			SetHeader("Content-Type", "application/octet-stream").
			SetBody(string(fileContents)).
			Send(c)

	default:
		response.SetStatusCode(404).Send(c)
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
