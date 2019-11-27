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
)

// List processes from which we need approval
var pWait list.List
var pDiff list.List
var criticalSection bool

func demandWait(ch chan bool) {
	for pWait.Len() != 0 {}
	criticalSection = true
	ch <- true
}

// Main entrypoint for the mutual exclusion program
func main() {
	// Create channels which communicate with the Client Process
	demand := make(chan bool)
	wait := make(chan bool)
	end := make(chan bool)

	var timestamp uint16
	//var demandTimestamp uint16
	//var currentDemand bool

	var parameters network.Parameters

	// Launch Client Process
	go client.PromptClient(demand, wait, end)

	// Infinite loop
	for {
		select {
		case <-demand:
			timestamp++
	//		currentDemand = true
	//		demandTimestamp = timestamp

			for i := uint8(0); i < parameters.NbProcesses ; i++  {
				// TODO REQ(currentDemand,i)
			}

			go demandWait(wait)

		case <-end:
			timestamp++
			criticalSection = false
	//		currentDemand = false

			for e := pDiff.Front(); e != nil; e = e.Next() {
				// TODO OK(timestamp, e.Value)
			}

			pWait = pDiff
			pDiff.Init()
		}
	}
}
