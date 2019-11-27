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
