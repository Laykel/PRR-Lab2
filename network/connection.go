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

var (
	entering = make(chan chan<- string) // Channel of channels
	leaving  = make(chan chan<- string)
	messages = make(chan string)
)

// Main TCP server entrypoint
func Server(req chan RequestCS, ok chan ReleaseCS, processNbr int) {
    // Listen on the current process' port
    recipientPort := strconv.Itoa(int(Params.InitialPort + uint16(processNbr)))
	listener, err := net.Listen("tcp", "127.0.0.1:" + recipientPort)
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
		go handleConn(conn, req)
	}
}

// Manage a specific TCP connection
func handleConn(conn net.Conn, req chan RequestCS) {
    defer conn.Close()

	ch := make(chan string)
	go func() {
		for msg := range ch {
			fmt.Println(msg)
		}
	}()

	who := conn.RemoteAddr().String()
	ch <- "You are " + who           // clientwriter <- handleConn
	fmt.Println(conn.RemoteAddr().String())
	req <- RequestCS{1, 23, 2}

	input := bufio.NewScanner(conn)
	for input.Scan() { // handleConn <- netcat client
		messages <- who + ": " + input.Text() // broadcaster <- handleConn
	}

	leaving <- ch
	messages <- who + " has left" // broadcaster <- handleConn
}

// Send bytes to recipient (port number calculated from initial port)
func Send(message []byte, recipient uint8) {
    // Connect to recipient's server
    recipientPort := strconv.Itoa(int(Params.InitialPort + uint16(recipient)))
    conn, err := net.Dial("tcp", "127.0.0.1:" + recipientPort)
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    // Send encoded message
    _, err = fmt.Fprintln(conn, message)
    if err != nil {
        log.Fatal(err)
    }
}
