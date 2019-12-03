/*
Lab 2 - mutual exclusion
File: mutex/mutex.go
Authors: Jael Dubey, Luc Wachter
Go version: 1.13.4 (linux/amd64)

Main entrypoint for the mutual exclusion program.

The access to the shared variable is guaranteed to be mutually exclusive
thanks to the Carvalho-Roucairol algorithm.

This file contains the central part of the algorithm, receiving requests
from the client and forwarding them to the network manager.
*/
package main

import (
    "container/list"
    "../client"
    "../network"
    "os"
    "strconv"
)

// TODO try and reorganize variables and functions (file for main and file for mutex?)
// List processes from which we need approval
var pWait list.List
var pDiff list.List

var timestamp uint32
var demandTimestamp uint16
var currentDemand bool

type listElem struct {
	ProcessNbr uint8
}

var criticalSection bool

func demandWait(ch chan bool) {
	// TODO check if this is legal in every Canton
	for pWait.Len() != 0 {
	}
	criticalSection = true
	ch <- true
}

func makeDemand(idArg int, wait chan bool) {
	timestamp++
	currentDemand = true
	demandTimestamp = uint16(timestamp)

	// For every process in pWait
	for e := pWait.Front(); e != nil; e = e.Next() {
		request := network.RequestCS{
			ReqType:    network.ReqType,
			ProcessNbr: uint8(idArg),
			Timestamp:  timestamp,
		}

		// Encode message and send to recipient
		network.Send(network.Encode(request), e.Value.(listElem).ProcessNbr)
	}

	demandWait(wait)
}

func endDemand(idArg int, val int32) {
	timestamp++
	criticalSection = false
	currentDemand = false

	// For every process in pDiff
	for e := pDiff.Front(); e != nil; e = e.Next() {
		ok := network.ReleaseCS{
			ReqType:    network.OkType,
			ProcessNbr: uint8(idArg),
			Timestamp:  timestamp,
			Value:      val,
		}

		// Encode message and send to recipient
		network.Send(network.Encode(ok), e.Value.(listElem).ProcessNbr)
	}

	pWait = pDiff
	pDiff.Init()
}

// Main entrypoint for the mutual exclusion program
func main() {
	// Create channels to communicate with the Client Process
	demand := make(chan bool)
	wait := make(chan bool)
	end := make(chan int32) // int32 to pass the new value
	quit := make(chan bool)

	// Create channels to communicate with the Network Process
	req := make(chan network.RequestCS)
	ok := make(chan network.ReleaseCS)

    // Handle program arguments
    // TODO Possible to move in client??? Or in main.go
    var idArg int

    if len(os.Args) == 2 {
        idArg, _ = strconv.Atoi(os.Args[1])
    } else {
        idArg = 0
    }

    // Launch Server Process
    go network.Server(req, ok, idArg)

	// Launch Client Process
	go client.PromptClient(demand, wait, end, quit)

	// Infinite loop
	for {
		select {
		// Client asks for critical section
		case <-demand:
			go makeDemand(idArg, wait)

        // Client releases critical section
		case val := <-end:
			go endDemand(idArg, val)

        case <-quit:
            return
		}
	}
}
