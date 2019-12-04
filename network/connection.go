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
    "log"
    "net"
)

// Main TCP server entrypoint
func Listen(address string, req chan []byte) {
	// Listen for incoming traffic
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print("Error accepting connections: ", err)
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

	// Send byte array to mutex
	req <- message
}

// Send bytes to recipient
func Send(message []byte, address string) {
	// Connect to recipient's server
	conn, err := net.Dial("tcp", address)
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
