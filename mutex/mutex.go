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
    "../client"
    "../network"
    "os"
	"strconv"
)

// TODO try and reorganize variables and functions (file for main and file for mutex?)
// List processes from which we need approval
//var pWait list.List
var pWait = make(map[uint8]bool)
var pDiff = make(map[uint8]bool)

var timestamp uint32
var demandTimestamp uint32
var currentDemand bool

type listElem struct {
	ProcessNbr uint8
}

var criticalSection bool

// Max returns the larger of x or y.
func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func demandWait(ch chan bool) {
	// TODO check if this is legal in every Canton
	for len(pWait) != 0 {
	}
	criticalSection = true
	ch <- true
}

func makeDemand(idArg int, wait chan bool) {
	timestamp++
	currentDemand = true
	demandTimestamp = timestamp

	// For every process in pWait
	for k := range pWait {
		request := network.RequestCS{
			ReqType:    network.ReqType,
			ProcessNbr: uint8(idArg),
			Timestamp:  timestamp,
		}

		// Encode message and send to recipient
		network.Send(network.Encode(request), k)
	}

	demandWait(wait)
}

func endDemand(idArg int, val int32) {
	timestamp++
	criticalSection = false
	currentDemand = false

	for k := range pDiff {
		ok := network.ReleaseCS{
			ReqType:    network.OkType,
			ProcessNbr: uint8(idArg),
			Timestamp:  timestamp,
			Value:      val,
		}

		// Encode message and send to recipient
		network.Send(network.Encode(ok), k)
	}

	pWait = pDiff
	pDiff = make(map[uint8]bool)
}

func okReceive(receiveTimestamp uint32, processNbr uint8) {
	timestamp = uint32(Max(int64(timestamp), int64(receiveTimestamp)) + 1)

	delete(pWait, processNbr)
}

func recReceive(idArg int, val int32, receiveTimestamp uint32, processNbr uint8) {
	timestamp = uint32(Max(int64(timestamp), int64(receiveTimestamp)) + 1)

	if currentDemand == false {
		ok := network.ReleaseCS{
			ReqType:    network.OkType,
			ProcessNbr: uint8(idArg),
			Timestamp:  timestamp,
			Value:      val,
		}

		// Encode message and send to recipient
		network.Send(network.Encode(ok), processNbr)
		pWait[processNbr] = true
	} else {
		if criticalSection || demandTimestamp < receiveTimestamp || ((demandTimestamp == receiveTimestamp) && (idArg < int(processNbr))) {
			pDiff[processNbr] = true
		} else {
			ok := network.ReleaseCS{
				ReqType:    network.OkType,
				ProcessNbr: uint8(idArg),
				Timestamp:  timestamp,
				Value:      val,
			}

			// Encode message and send to recipient
			network.Send(network.Encode(ok), processNbr)
			pWait[processNbr] = true

			request := network.RequestCS{
				ReqType:    network.ReqType,
				ProcessNbr: uint8(idArg),
				Timestamp:  timestamp,
			}

			// Encode message and send to recipient
			network.Send(network.Encode(request), processNbr)

		}
	}
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

		// Other site releases critical section via Network
		case okReceived := <-ok:
			go okReceive(okReceived.Timestamp, okReceived.ProcessNbr)

		// Other site asks for critical section
		case reqReceived := <- req:
			go recReceive(idArg, 0, reqReceived.Timestamp, reqReceived.ProcessNbr)

        case <-quit:
            return
		}
	}
}
