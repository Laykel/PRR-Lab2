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
)

var parameters Parameters

var (
	entering = make(chan chan<- string) // Channel of channels
	leaving  = make(chan chan<- string)
	messages = make(chan string)
)

func Server() {
	listener, err := net.Listen("tcp", "localhost:" + string(parameters.InitialPort))
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[chan<- string]bool) // all connected clients
	for {
		select {
		case msg := <-messages: // broadcaster <- handleConn
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli <- msg // clientwriter (handleConn) <- broadcaster
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string) // channel 'client' mais utilisé ici dans les 2 sens
	go func() {             // clientwriter
		for msg := range ch { // clientwriter <- broadcaster, handleConn
			fmt.Fprintln(conn, msg) // netcat Client <- clientwriter
		}
	}()

	who := conn.RemoteAddr().String()
	ch <- "You are " + who           // clientwriter <- handleConn
	messages <- who + " has arrived" // broadcaster <- handleConn
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan() { // handleConn <- netcat client
		messages <- who + ": " + input.Text() // broadcaster <- handleConn
	}

	leaving <- ch
	messages <- who + " has left" // broadcaster <- handleConn
	conn.Close()
}
