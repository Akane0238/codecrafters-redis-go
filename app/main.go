package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func handle(conn net.Conn, myRedis *KVStore) {
	defer conn.Close()

	buffer := make([]byte, 128)
	for {
		n, err := conn.Read(buffer)
		if n > 0 {
			var parser RESPParser

			// buffer --> parser --> client input
			var elements []string = parser.Decode(string(buffer))

			// msg --> parser --> server response
			switch strings.ToUpper(elements[0]) {
			case "ECHO":
				response := parser.Encode(elements[1])
				conn.Write([]byte(response))

			case "PING":
				conn.Write([]byte("+PONG\r\n"))

			case "SET":
				myRedis.Set(elements[1], elements[2])
				conn.Write([]byte("+OK\r\n"))

			case "GET":
				val, ok := myRedis.Get(elements[1])
				if ok {
					res := fmt.Sprintf("$%d\r\n%s\r\n", len(val), val)
					conn.Write([]byte(res))
				} else {
					conn.Write([]byte("$-1\r\n")) // null bulk string
				}

			default:
				fmt.Println("Unsupported command: ", elements[0])

			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			fmt.Println("Error reading from client: ", err.Error())
			break
		}

	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	// 初始化共享数据库
	myRedis := NewKVStore()

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
		go handle(conn, myRedis)
	}
}
