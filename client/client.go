/*
Lab 2 - mutual exclusion
File: client/client.go
Authors: Jael Dubey, Luc Wachter
Go version: 1.13.4 (linux/amd64)

Provides a simple console user interface to access and modify a variable
Shared between all processes.
*/
package client

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Shared variable across processes
var Shared int32

// Display a prompt for the user with instructions
func prompt() {
	fmt.Println("Commands: [r to read variable], [w <integer> to write to variable], [q to quit].")
	fmt.Print("> ")
}

// TODO errors checking
// Ask user for their choice and either prints value or ask for CS and modify value
func PromptClient(demand chan bool, wait chan bool, end chan int32, quit chan bool) {

	reader := bufio.NewReader(os.Stdin)

	for {
		// Ask the user what he wants to do
		prompt()
		input, _ := reader.ReadString('\n')

		tokens := strings.Split(input[:len(input)-1], " ")

		switch tokens[0] {
		// The user wants to read the variable
		case "r":
			// Just print the variable to stdout
			fmt.Println(Shared)

		// The user wants to write to the variable
		case "w":
			// TODO check value
			newValue, _ := strconv.ParseInt(tokens[1], 10, 32)
			fmt.Println(Shared)

			// Call the Carvalho - Roucairol algorithm to acquire critical section
			demand <- true
			// Wait until the CS is free
			<-wait
			// START of critical section

			// Modify the variable
			Shared = int32(newValue)

			// END of critical section
			// Then liberate the critical section
			end <- Shared

			fmt.Println(Shared)

		// The user wants to quit the program
		case "q":
			quit <- true
			return

		// Unknown command
		default:
			fmt.Println("This is not one of the allowed commands. Please read the instructions.")
		}
	}
}
