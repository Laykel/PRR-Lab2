/*
Lab 2 - mutual exclusion
File: client/client.go
Authors: Jael Dubey, Luc Wachter
Go version: 1.13.4 (linux/amd64)

Provides a simple console user interface to access and modify a variable
shared between all processes.

The access to the shared variable is guaranteed to be mutually exclusive
thanks to the Carvalho-Roucairol algorithm.
 */
package client

import (
    "../network"
    "bufio"
    "fmt"
    "os"
    "strconv"
)

// Ask user for their choice and either prints value or ask for CS and modify value
func clientProcess(commWithMutexProcess chan int) {
    // Shared variable across processes
    var shared int32
    reader := bufio.NewReader(os.Stdin)
    input,_ := reader.ReadString('\n')
    // Ask the user what he wants to do
    // Allow him to read or write the shared variable

    // TODO: For loop
    // CASE READ
    if input[0] == 'r' {
        // Just prints the variable to stdout
        fmt.Println(shared)
    }

    if input[0] == 'w' {
        // CASE WRITE
        newValue, _ := strconv.Atoi(input[1:])
        // Prints the variable to stdout
        fmt.Println(shared)
        // Calls the Carvalho - Roucairol algorithm to acquire critical section
        commWithMutexProcess <- network.Demand
        commWithMutexProcess <- network.Wait
        // Then modifies the variable
        shared = int32(newValue)
        // Then notifies the other processes
        commWithMutexProcess <- network.End
        // Then liberates the critical section
    }
}
