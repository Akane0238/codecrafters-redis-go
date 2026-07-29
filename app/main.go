package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func handle(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 128)
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading from client: ", err.Error())
			break
		}

		var parser RESPParser

		// buffer --> parser --> client input
		var elements []string = parser.Decode(string(buffer))

		// msg --> parser --> server response
		if elements[0] == "ECHO" {
			response := parser.Encode(elements[1])
			conn.Write([]byte(response))
		}
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			continue
		}

		// Receive PING and reponse in hardcode
		go handle(conn)
	}
}
