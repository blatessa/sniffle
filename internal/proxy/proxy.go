package proxy

import (
	"errors"
	"net"

	"github.com/blatessa/sniffle/internal/request"
)

type Proxy struct{}

func (p *Proxy) Forward(request *request.Request) (string, error) {
	host := request.Headers["Host"]
	if host == "" {
		return "", errors.New("missing Host header")
	}
	port := "80"
	if request.Headers["Port"] != "" {
		port = request.Headers["Port"]
	}

	conn, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		return "", errors.New("failed to connect to server")
	}

	defer conn.Close()

	conn.Write([]byte(request.Method + " " + request.Path + " HTTP/1.1\r\n"))
	for key, value := range request.Headers {
		conn.Write([]byte(key + ": " + value + "\r\n"))
	}
	conn.Write([]byte("\r\n"))

	responseBuffer := make([]byte, 4096)
	n, err := conn.Read(responseBuffer)
	if err != nil {
		return "", errors.New("failed to read response from server")
	}

	response := string(responseBuffer[:n])
	if response == "" {
		return "", errors.New("empty response from server")
	}

	return response, nil
}
