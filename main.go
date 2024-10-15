package main

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Client connected:", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	for {
		message, _ := reader.ReadString('\n')
		fmt.Println("Message received:", message)
		conn.Write([]byte("Message received: " + message))
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server listening on port 8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}

		go func() {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("PANIC", slog.Any("panic", r))
				}

				fmt.Println("Client disconnected:", conn.RemoteAddr())
				conn.Close()
			}()

			handleConnection(conn) // Spawns a new goroutine for each connection
		}()
	}
}
