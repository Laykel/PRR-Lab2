/*
Main entrypoint for the mutual exclusion test application.

Provides a simple console user interface to access and modify a variable
shared between all processes.

The access to the shared variable is guaranteed to be mutually exclusive
thanks to the Carvalho-Roucairol algorithm.
 */
package main

import "fmt"

// TODO: Maybe, this should be a go routine triggered by a controller
func main() {
    var shared int64

    // Ask the user what he wants to do
    // Allow him to read or write the shared variable

    // CASE READ
    // Just prints the variable to stdout
    fmt.Println(shared)

    // CASE WRITE
    // Prints the variable to stdout
    // Calls the Carvalho - Roucairol algorithm to acquire critical section
    // Then modifies the variable
    // Then notifies the other processes
    // Then liberates the critical section
}
