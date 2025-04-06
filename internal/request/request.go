package request

import (
	"errors"
	"net/url"
	"strings"
)

type Request struct {
	Method  string
	Host    string
	Port    string
	Path    string
	Headers map[string]string
}

func Parse(data []byte) (*Request, error) {
	request := new(Request)
	lines := strings.Split(string(data), "\r\n")
	if len(lines) < 1 {
		return nil, errors.New("invalid request: no request line")
	}

	requestLine := lines[0]
	headerLines := lines[1:]

	method, path, err := parseStartLine(requestLine)
	if err != nil {
		return nil, err
	}

	headers := parseHeaders(headerLines)

	host, ok := headers["Host"]
	if !ok {
		return nil, errors.New("missing Host header")
	}

	port := "80"
	if p, ok := headers["Port"]; ok {
		port = p
	}

	request.Method = method
	request.Path = path
	request.Host = host
	request.Port = port
	request.Headers = headers

	return request, nil
}

func parseStartLine(requestLine string) (string, string, error) {
	parts := strings.SplitN(requestLine, " ", 3)
	if len(parts) < 2 {
		return "", "", errors.New("malformed request line")
	}

	method := parts[0]
	uri := parts[1]

	u, err := url.Parse(uri)
	if err != nil {
		return "", "", errors.New("unable to parse URL")
	}

	return method, u.RequestURI(), nil
}

func parseHeaders(lines []string) map[string]string {
	headers := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}
	return headers
}
