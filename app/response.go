package main

import (
	"fmt"
	"net"
	"net/http"
)

var statusCodes = map[int]string{
	200: "OK",
	201: "Created",
	404: "Not Found",
}

type Response struct {
	Version    string
	StatusCode int
	Headers    map[string]string
	Content    string
}

func (r *Response) AddStatus(code int) *Response {
	r.StatusCode = code
	return r
}

func (r *Response) AddHeader(key string, value string) *Response {
	r.Headers[key] = value
	return r
}

func (r *Response) AddContent(data string) *Response {
	r.Content = data
	r.AddHeader("Content-Length", fmt.Sprint(len(data)))
	return r
}

func (r *Response) Write(conn net.Conn) {
	respStr := fmt.Sprintf("HTTP/1.1 %d %s\r\n", r.StatusCode, statusCodes[r.StatusCode])

	for key, val := range r.Headers {
		respStr += (key + ": " + val + "\r\n")
	}

	respStr += ("\r\n" + r.Content + "\r\n")

	conn.Write(
		[]byte(respStr),
	)
}

func NewResponse() *Response {
	return &Response{
		Version:    "HTTP/1.1",
		StatusCode: http.StatusOK,
		Headers: map[string]string{
			"Content-Type":   "text/plain",
			"Content-Length": "0",
		},
		Content: "",
	}
}
