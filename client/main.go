package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"sync"
)

func handleError(err error, message string) {
	if err != nil {
		fmt.Printf("Error: %s %s", message, err.Error())
		os.Exit(1)
	}
}

func sendMessageToServer(conn net.Conn, wg *sync.WaitGroup) {
    defer conn.Close()
    defer wg.Done()
	fmt.Printf("Enter your name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')
	if len(name) > 0 {
		name = name[:len(name)-1]
	}

	for {
		text, _ := reader.ReadString('\n')
		if text == "exit\n" {
			_, err := conn.Write([]byte("exit\n"))
			handleError(err, "Error sending message to server")
			break
		} else {
			finalMessage := fmt.Sprintf("%s: %s", name, text)
			_, err := conn.Write([]byte(finalMessage))
			handleError(err, "Error sending message to server")
		}
	}
}

func readMessageFromServer(conn net.Conn, wg *sync.WaitGroup){
    defer conn.Close()
    defer wg.Done()
    reader := bufio.NewReader(conn)
    for {
        message, err := reader.ReadString('\n')
        handleError(err, "Error reading message from server")
        fmt.Printf("%s", message)
    }
}

func main() {
    // creating a connection with the server
	conn, err := net.Dial("tcp", "localhost:2000")
	handleError(err, "Error connecting to server")

    var wg sync.WaitGroup
    wg.Add(2)

    // sending message to the server
	go sendMessageToServer(conn, &wg)

    // reading message from the server
    go readMessageFromServer(conn, &wg)
    wg.Wait()
}
