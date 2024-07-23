package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	// "strings"
)

func handleError(err error, message string) bool {
	if err != nil {
		if err == io.EOF {
			return true
		}
		fmt.Printf("Error: %s %s", message, err.Error())
		os.Exit(1)
	}
	return false
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
        closeConnection := handleError(err, "Error reading:")

        // closing connection if client closes the connection
        if closeConnection {
            break
        }

		fmt.Print(string(message))
        
		if string(message) == "exit\n" {
			fmt.Printf("breaking connection\n")
			break
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:2000")
	handleError(err, "Error listening:")
	defer listener.Close()

	fmt.Println("Server is listening...")

	for {
		conn, err := listener.Accept()
		handleError(err, "Error accepting:")
		go handleRequest(conn)
	}
}
