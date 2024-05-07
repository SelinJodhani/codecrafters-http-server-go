package main

import (
	"fmt"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Body    string
}

func ParseHTTPRequest(request string) (*Request, error) {
	lines := strings.Split(request, "\r\n")
	headers := make(map[string]string)

	var (
		method  string
		path    string
		version string
		body    string
	)

	specs := strings.Fields(lines[0])

	if len(specs) != 3 {
		return nil, fmt.Errorf("malformed request")
	}

	method = specs[0]
	path = specs[1]
	version = specs[2]

	for i := 1; i < len(lines); i++ {
		line := lines[i]

		if line == "" {
			// Empty line indicates end of headers
			body = strings.Join(lines[i+1:], "\r\n")
			break
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			headers[key] = value
		}
	}

	return &Request{
		Method:  method,
		Path:    path,
		Version: version,
		Headers: headers,
		Body:    body,
	}, nil
}
