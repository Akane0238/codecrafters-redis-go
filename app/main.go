package main

import (
	"fmt"
	"net"
	"os"
)

// Ensures gofmt doesn't remove the "net" and "os" imports in stage 1 (feel free to remove this!)
var _ = net.Listen
var _ = os.Exit

func handle(conn net.Conn){
	defer conn.Close()

	buffer := make([]byte, 128)
	for {
		n, err := conn.Read(buffer)
		if err != nil{
			fmt.Println("Error reading from client: ", err.Error)
			os.Exit(1)
		}
		fmt.Println("Receive from client: ", string(buffer[:n]))

		conn.Write([]byte("+PONG\r\n"))
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	listener, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	// Receive PING and reponse in hardcode
	go handle(conn)

}
