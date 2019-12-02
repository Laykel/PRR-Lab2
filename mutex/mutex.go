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
	"container/list"
	"os"
	"strconv"
)

// TODO try and reorganize variables and functions (file for main and file for mutex?)
// List processes from which we need approval
var pWait list.List
var pDiff list.List

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

// Main entrypoint for the mutual exclusion program
func main() {
	// Create channels which communicate with the Client Process
	demand := make(chan bool)
	wait := make(chan bool)
	end := make(chan int32)

	// Create channels which communicate with the Network Process
	req := make(chan network.RequestCS)
	ok := make(chan network.ReleaseCS)

	var timestamp uint32
	//var demandTimestamp uint16
	//var currentDemand bool

	var parameters network.Parameters

	// Launch Client Process
	go client.PromptClient(demand, wait, end)

	// Launch Server Process
	go network.Server(req, ok)

	// TODO Handle arguments
	tmp, _ := strconv.Atoi(os.Args[1])

	// Infinite loop
	for {
		select {
		case <-demand:
			timestamp++
			//		currentDemand = true
			//		demandTimestamp = timestamp

			for e := pWait.Front(); e != nil; e = e.Next() {
				request := network.RequestCS{
					ReqType:    network.REQ_TYPE,
					ProcessNbr: uint8(tmp),
					Timestamp:  timestamp,
				}

				network.SendReq(request, e.Value.(listElem).ProcessNbr)
			}

			go demandWait(wait)

		case val := <-end:
			timestamp++
			criticalSection = false
			//		currentDemand = false

			for e := pDiff.Front(); e != nil; e = e.Next() {
				ok := network.ReleaseCS{
					ReqType:    network.OK_TYPE,
					ProcessNbr: uint8(tmp),
					Timestamp:  timestamp,
					Value:      val,
				}

				network.SendOk(ok, e.Value.(listElem).ProcessNbr)
			}

			pWait = pDiff
			pDiff.Init()
		}
	}
}
