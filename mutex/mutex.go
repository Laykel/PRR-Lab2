/*
Lab 2 - mutual exclusion
File: mutex/mutex.go
Authors: Jael Dubey, Luc Wachter
Go version: 1.13.4 (linux/amd64)

This file contains the implementation of the Carvalho-Roucairol algorithm.
*/
package main

import (
	"github.com/Laykel/PRR-Lab2/network"
)

// List processes from which we need approval
var pWait = make(map[uint8]bool)
var pDiff = make(map[uint8]bool)

// Necessary algorithm variables
var timestamp uint32
var demandTimestamp uint32
var currentDemand bool

var criticalSection bool

// Max returns the larger of x and y
func Max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}

func makeDemand(processId uint8, wait chan bool) {
	timestamp++
	currentDemand = true
	demandTimestamp = timestamp

	// For every process in pWait
	for k := range pWait {
		request := network.RequestCS{
			ReqType:    network.ReqType,
			ProcessNbr: processId,
			Timestamp:  timestamp,
		}

		// Encode message and send to recipient
		network.Send(network.Encode(request), k)
	}

	// TODO check if this is legal in every Canton
	for len(pWait) != 0 {
	}
	criticalSection = true
	wait <- true
}

func endDemand(processId uint8, val int32) {
	timestamp++
	criticalSection = false
	currentDemand = false

	for k := range pDiff {
		ok := network.ReleaseCS{
			ReqType:    network.OkType,
			ProcessNbr: processId,
			Timestamp:  timestamp,
		}

		// Encode message and send to recipient
		network.Send(network.Encode(ok), k)
	}

	pWait = pDiff
	pDiff = make(map[uint8]bool)
}

func okReceive(ok network.ReleaseCS) {
	timestamp = uint32(Max(int64(timestamp), int64(ok.Timestamp)) + 1)

	delete(pWait, ok.ProcessNbr)
}

func reqReceive(processId uint8, req network.RequestCS) {
	timestamp = uint32(Max(int64(timestamp), int64(req.Timestamp)) + 1)

	if currentDemand == false {
		ok := network.ReleaseCS{
			ReqType:    network.OkType,
			ProcessNbr: processId,
			Timestamp:  timestamp,
		}

		// Encode message and send to recipient
		network.Send(network.Encode(ok), req.ProcessNbr)
		pWait[req.ProcessNbr] = true
	} else {
		if criticalSection || demandTimestamp < req.Timestamp || ((demandTimestamp == req.Timestamp) && (processId < req.ProcessNbr)) {
			pDiff[req.ProcessNbr] = true
		} else {
			ok := network.ReleaseCS{
				ReqType:    network.OkType,
				ProcessNbr: processId,
				Timestamp:  timestamp,
			}

			// Encode message and send to recipient
			network.Send(network.Encode(ok), req.ProcessNbr)
			pWait[req.ProcessNbr] = true

			request := network.RequestCS{
				ReqType:    network.ReqType,
				ProcessNbr: processId,
				Timestamp:  timestamp,
			}

			// Encode message and send to recipient
			network.Send(network.Encode(request), req.ProcessNbr)
		}
	}
}
