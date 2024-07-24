package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

type connObj struct {
    conn net.Conn;
    address string;
}

var connections []connObj


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

		if closeConnection ||  string(message) == "exit\n" {
            var tempConnections []connObj
            for _, connection := range connections {
                if connection.address != conn.RemoteAddr().String() {
                    tempConnections = append(tempConnections, connection)
                }
            }
            connections = tempConnections
			break
		}

        // handle this in its own go routine
        for _, connection := range connections {
            if connection.address != conn.RemoteAddr().String() {
                _, err := connection.conn.Write([]byte(string(message)))
                handleError(err, "Error while sending messges to other servers\n")
            }
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
        connections = append(connections, connObj{conn: conn, address: conn.RemoteAddr().String()})
		go handleRequest(conn)
	}
}
