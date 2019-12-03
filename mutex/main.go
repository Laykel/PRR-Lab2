/*
Lab 2 - mutual exclusion
File: mutex/main.go
Authors: Jael Dubey, Luc Wachter
Go version: 1.13.4 (linux/amd64)

Main entrypoint for the mutual exclusion program.

The access to the shared variable is guaranteed to be mutually exclusive
thanks to the Carvalho-Roucairol algorithm.

This file contains the central part of the algorithm, receiving requests
from the client, forwarding them to the network manager and calling implementation functions.
*/
package main

import (
	"encoding/json"
	"fmt"
	"github.com/Laykel/PRR-Lab2/client"
	"github.com/Laykel/PRR-Lab2/network"
	"os"
	"strconv"
)

// Path to json parameters file
const parametersFile = "mutex/parameters.json"

// Load parameters from json file
func loadParameters(file string) network.Parameters {
	var params network.Parameters

	// Read parameters file
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&params)

	return params
}

// Main entrypoint for the mutual exclusion program
func main() {
	// Create channels to communicate with the Client Process
	demand := make(chan bool)
	wait := make(chan bool)
	end := make(chan int32) // int32 to pass the new value
	quit := make(chan bool)

	// Create channel to get message from network process
	message := make(chan []byte)

	// Handle program arguments
	var processId uint8

	if len(os.Args) == 2 {
		tmp, _ := strconv.Atoi(os.Args[1])
		processId = uint8(tmp)
	} else {
		processId = 0
	}

	// Read parameters from json file
	network.Params = loadParameters(parametersFile)

	// Populate pWait
	for i := uint8(0); i < network.Params.NbProcesses; i++ {
		if i != processId {
			pWait[i] = true
		}
	}

	// Launch Server Process
	go network.Listen(processId, message)

	// Launch Client Process
	go client.PromptClient(demand, wait, end, quit)

	// Infinite loop
	for {
		select {
		// Client asks for critical section
		case <-demand:
			go makeDemand(processId, wait)

			// Client releases critical section
		case val := <-end:
		    // TODO send setValue
			go endDemand(processId, val)

		// Other site releases critical section via Network
		case receivedMsg := <-message:
			// Handle request type
			switch receivedMsg[0] {
			case byte(network.ReqType):
				req := network.DecodeRequest(receivedMsg)
				go reqReceive(processId, req)

			case byte(network.OkType):
				ok := network.DecodeRelease(receivedMsg)
				go okReceive(ok)

			case byte(network.ValType):
                // TODO
			}

		case <-quit:
			return
		}
	}
}