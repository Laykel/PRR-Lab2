/*
Lab 2 - mutual exclusion
File: network/connection.go
Authors: Jael Dubey, Luc Wachter
Go version: 1.13.4 (linux/amd64)

Handle TCP connections, forward requests to and from the mutex controller

Source: https://go-talks.appspot.com/github.com/patricklac/prr-slides/ch2
*/
package network

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
)

// Main TCP server entrypoint
func Listen(processNbr uint8, req chan []byte) {
	// Listen on the current process' port
	recipientPort := strconv.Itoa(int(Params.InitialPort + uint16(processNbr)))
	listener, err := net.Listen("tcp", "127.0.0.1:"+recipientPort)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		// Manage this connection without blocking so that we don't miss connections
		go handleConnection(conn, req)
	}
}

// Manage a specific TCP connection
func handleConnection(conn net.Conn, req chan []byte) {
	defer conn.Close()

    // Read from conn
	input := bufio.NewScanner(conn)
	input.Scan()

	message := input.Bytes()

	fmt.Println(conn.RemoteAddr().String())

	// Send byte array to mutex
	req <- message
}

// Send bytes to recipient (port number calculated from initial port)
func Send(message []byte, recipient uint8) {
	// Connect to recipient's server
	recipientPort := strconv.Itoa(int(Params.InitialPort + uint16(recipient)))
	conn, err := net.Dial("tcp", "127.0.0.1:"+recipientPort)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Send encoded message
    _, err = conn.Write(message)
	if err != nil {
		log.Fatal(err)
	}
}
