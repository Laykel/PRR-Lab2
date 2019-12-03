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
func Server(req chan RequestCS, ok chan ReleaseCS) {
	listener, err := net.Listen("tcp", "localhost:"+string(Params.InitialPort))
	if err != nil {
		log.Fatal(err)
	}

	go broadcaster(req)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster(req chan RequestCS) {
	clients := make(map[chan<- string]bool) // all connected clients
	for {
		select {
		case msg := <-messages: // broadcaster <- handleConn
			// Broadcast incoming message to all clients' outgoing message channels.
			for cli := range clients {
				cli <- msg // clientwriter (handleConn) <- broadcaster
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli)

		/*case request := <-req:
			reqType := request.ReqType
			reqTimestamp := request.Timestamp
			reqProcessNb := request.ProcessNbr

			if request.ReqType == REQ_TYPE {
				for cli := range clients {

				}
			}*/


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

// Send bytes to recipient (port number calculated from initial port)
func Send(message []byte, recipient uint8) {
    conn, err := net.Dial("tcp", "localhost:" + strconv.Itoa(int(Params.InitialPort + uint16(recipient))))
    if err != nil {
        log.Fatal(err)
    }

    // Send encoded message
    _, err = fmt.Fprint(conn, message)
    if err != nil {
        log.Fatal(err)
    }

    conn.Close()
}
