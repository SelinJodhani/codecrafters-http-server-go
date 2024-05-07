package main

import (
	"fmt"
	"net"
	"net/http"
)

type Response struct {
	Version       string
	StatusCode    int
	ContentType   string
	ContentLength int
	Content       string
}

func (r *Response) AddStatus(code int) *Response {
	r.StatusCode = code
	return r
}

func (r *Response) AddContentType(contentType string) *Response {
	r.ContentType = contentType
	return r
}

func (r *Response) AddContent(data string) *Response {
	r.Content = data
	r.ContentLength = len(data)
	return r
}

func (r *Response) Write(conn net.Conn) {
	conn.Write(
		[]byte(
			fmt.Sprintf(
				"HTTP/1.1 %d OK\r\nContent-Type: %s\r\nContent-Length: %d\r\n\n%s\r\n",
				r.StatusCode,
				r.ContentType,
				r.ContentLength,
				r.Content,
			),
		),
	)
}

func NewResponse() *Response {
	return &Response{
		Version:       "HTTP/1.1",
		StatusCode:    http.StatusOK,
		ContentType:   "text/plain",
		ContentLength: 0,
		Content:       "",
	}
}
