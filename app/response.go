package main

import (
	"bytes"
	"compress/gzip"
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

func (r *Response) SetStatusCode(code int) *Response {
	r.StatusCode = code
	return r
}

func (r *Response) SetHeader(key string, value string) *Response {
	r.Headers[key] = value
	return r
}

func (r *Response) SetBody(body string) *Response {
	if encoding, ok := r.Headers["Content-Encoding"]; ok && encoding == "gzip" {
		data := []byte(body)

		var compressed bytes.Buffer
		gzipWriter := gzip.NewWriter(&compressed)

		_, err := gzipWriter.Write(data)
		if err != nil {
			fmt.Println("Error compressing data:", err)
			return r
		}

		err = gzipWriter.Close()
		if err != nil {
			fmt.Println("Error closing gzip writer:", err)
			return r
		}

		body = compressed.String()
	}

	r.Content = body
	r.SetHeader("Content-Length", fmt.Sprint(len(r.Content)))
	return r
}

func (r *Response) Send(conn net.Conn) {
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
