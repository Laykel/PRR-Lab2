/*
Lab 2 - mutual exclusion
File: client/client.go
Authors: Jael Dubey, Luc Wachter
Go version: 1.13.4 (linux/amd64)

Provides a simple console user interface to access and modify a variable
shared between all processes.
*/
package client

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// TODO comments, prompt, and errors checkin
// Ask user for their choice and either prints value or ask for CS and modify value
func PromptClient(demand chan bool, wait chan bool, end chan bool) {
	// Shared variable across processes
	var shared int32
	reader := bufio.NewReader(os.Stdin)

	// Ask the user what he wants to do
	// Allow him to read or write the shared variable

	for {
		input, _ := reader.ReadString('\n')

		tokens := strings.Split(input[:len(input)-1], " ")

		// CASE READ
		if tokens[0] == "r" {
			// Just prints the variable to stdout
			fmt.Println(shared)
		}

		if tokens[0] == "w" {
			// CASE WRITE
			newValue, _ := strconv.ParseInt(tokens[1], 10, 32)
			// Prints the variable to stdout
			fmt.Println(shared)
			// Calls the Carvalho - Roucairol algorithm to acquire critical section
			demand <- true
			// Wait until the CS is free
			<-wait
			// Then modifies the variable
			shared = int32(newValue)
			// Then notifies the other processes
			end <- true
			// Then liberates the critical section
		}
	}
}
