package proxy

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/blatessa/sniffle/internal/request"
	"github.com/blatessa/sniffle/internal/response"
)

type Proxy struct {
	requests  []*request.Request
	responses []*response.Response
}

func (p *Proxy) Forward(rawRequest []byte) ([]byte, error) {
	req, err := request.Parse(rawRequest)
	if err == nil {
		p.AddRequest(req)
	}

	conn, err := net.Dial("tcp", req.Host+":"+req.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to server: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write(rawRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	reader := bufio.NewReader(conn)

	var headerBuf bytes.Buffer
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, fmt.Errorf("failed to read response headers: %w", err)
		}
		headerBuf.WriteString(line)
		if line == "\r\n" {
			break
		}
	}

	headersRaw := headerBuf.String()
	contentLength := getContentLength(headersRaw)

	body := make([]byte, contentLength)
	_, err = io.ReadFull(reader, body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fullResp := append(headerBuf.Bytes(), body...)

	resp, err := response.Parse(fullResp)
	if err == nil {
		p.AddResponse(resp)
	}

	return fullResp, nil
}

func getContentLength(headers string) int {
	lines := strings.Split(headers, "\r\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.ToLower(line), "content-length:") {
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				val := strings.TrimSpace(parts[1])
				n, err := strconv.Atoi(val)
				if err == nil {
					return n
				}
			}
		}
	}
	return 0
}

func (p *Proxy) AddRequest(request *request.Request) {
	p.requests = append(p.requests, request)
}

func (p *Proxy) AddResponse(resp *response.Response) {
	p.responses = append(p.responses, resp)
}

func (p *Proxy) Dump() {
	fmt.Println("=== Dumping Proxy State ===")

	fmt.Println("\n-- Requests --")
	for i, req := range p.requests {
		fmt.Printf("[%d] %s %s %s:%s\n", i, req.Method, req.Path, req.Host, req.Port)
		for key, val := range req.Headers {
			fmt.Printf("    %s: %s\n", key, val)
		}
	}

	fmt.Println("\n-- Responses --")
	for i, res := range p.responses {
		fmt.Printf("[%d] %d\n", i, res.StatusCode)
		for key, val := range res.Headers {
			fmt.Printf("    %s: %s\n", key, val)
		}
		body := res.Body
		if len(body) > 100 {
			body = body[:100] + "..."
		}
		fmt.Println("    Body (truncated):")
		fmt.Println("    " + body)
	}
}
