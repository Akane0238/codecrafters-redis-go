package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

// Server 统一管理连接和依赖
type Server struct {
	db *KVStore // 全局共享数据库
	// ...
	// 可扩展全局资源
}

func (server *Server) handleConnection(conn net.Conn) {
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
				server.db.Set(elements[1], elements[2])
				conn.Write([]byte("+OK\r\n"))

			case "GET":
				val, ok := server.db.Get(elements[1])
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
