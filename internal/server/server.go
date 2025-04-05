package server

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/blatessa/sniffle/internal/request"
)

func Start(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	defer l.Close()

	fmt.Printf("Listening on %s...\n", address)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(c net.Conn) {
	defer c.Close()

	buf := make([]byte, 1024)
	n, err := c.Read(buf)
	if err != nil && err != io.EOF {
		log.Println("Error reading from connection:", err)
		return
	}
	if n == 0 {
		log.Println("Connection closed by client")
		return
	}

	parsed, err := request.Parse(buf)
	if err != nil {
		log.Println("Error parsing request:", err)
		return
	}
	fmt.Printf("%#v\n", parsed)
}
