package response

import (
	"errors"
	"strconv"
	"strings"
)

type Response struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

func Parse(data []byte) (*Response, error) {
	response := new(Response)

	parts := strings.SplitN(string(data), "\r\n\r\n", 2)
	if len(parts) < 2 {
		return nil, errors.New("invalid response: no header/body separator")
	}

	headerPart := parts[0]
	bodyPart := parts[1]
	lines := strings.Split(headerPart, "\r\n")

	statusCode, err := parseStatusLine(lines[0])
	if err != nil {
		return nil, err
	}

	headers, err := parseHeaders(lines[1:])
	if err != nil {
		return nil, err
	}

	response.StatusCode = statusCode
	response.Headers = headers
	response.Body = bodyPart

	return response, nil
}

func parseStatusLine(line string) (int, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) < 2 {
		return 0, errors.New("malformed status line")
	}

	statusCode, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, errors.New("invalid status code")
	}

	return statusCode, nil
}

func parseHeaders(lines []string) (map[string]string, error) {
	headers := make(map[string]string)
	for _, line := range lines {
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, errors.New("malformed header line")
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}
	return headers, nil
}
