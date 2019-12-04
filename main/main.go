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
	"../client"
	"../mutex"
	"../network"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Path to json parameters file
const parametersFile = "main/parameters.json"

// Load parameters from json file
func loadParameters(file string) network.Parameters {
	var params network.Parameters

	// Read parameters file
	configFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
	} else if configFile == nil {
	    fmt.Println("Could not open parameters file.")
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

	mutex.PopulatePWait(processId)

	// Launch Server Process
    port := strconv.Itoa(int(network.Params.InitialPort + uint16(processId)))
    address := network.Params.ProcessAddress+":"+port
	go network.Listen(address, message)

	// Launch Client Process
	go client.PromptClient(demand, wait, end, quit)

	for {
		select {
		// Client asks for critical section
		case <-demand:
			go mutex.MakeDemand(processId, wait)

        // Client releases critical section
		case val := <-end:
			go mutex.EndDemand(processId, val)

		// Other site releases critical section via Network
		case receivedMsg := <-message:
			switch receivedMsg[0] {
			case byte(network.RequestMessageType):
				req := network.DecodeRequest(receivedMsg)
				go mutex.ReqReceive(processId, req)

			case byte(network.ReleaseMessageType):
				ok := network.DecodeRelease(receivedMsg)
				go mutex.OkReceive(ok)

			case byte(network.SetValueMessageType):
				value := network.DecodeSetVariable(receivedMsg)
				client.Shared = value.Value
			}

		case <-quit:
			return
		}
	}
}
