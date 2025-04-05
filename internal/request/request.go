package request

import (
	"errors"
	"strings"
)

type Request struct {
	Method  string
	Path    string
	Headers map[string]string
}

func Parse(data []byte) (*Request, error) {
	request := new(Request)

	lines := strings.Split(string(data), "\r\n")

	if len(lines) < 1 {
		return nil, errors.New("invalid request: no request line")
	}

	requestLine := strings.TrimSpace(lines[0])
	headerLines := lines[1:]

	method, path := parseStartLine(requestLine)
	headerMap := parseHeaders(headerLines)

	request.Method = method
	request.Path = path
	request.Headers = headerMap

	return request, nil
}

func parseStartLine(requestLine string) (string, string) {

	parts := strings.SplitN(requestLine, " ", 3)

	method := parts[0]
	path := parts[1]

	return method, path
}

func parseHeaders(lines []string) map[string]string {
	headerMap := make(map[string]string)

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
		headerMap[key] = value
	}

	return headerMap
}
