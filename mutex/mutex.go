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
    "encoding/json"
    "fmt"
    "github.com/Laykel/PRR-Lab2/client"
	"github.com/Laykel/PRR-Lab2/network"
	"os"
	"strconv"
)

const parametersFile = "mutex/parameters.json"

// TODO try and reorganize variables and functions (file for main and file for mutex?)
// List processes from which we need approval
var pWait = make(map[uint8]bool)
var pDiff = make(map[uint8]bool)

var timestamp uint32
var demandTimestamp uint32
var currentDemand bool

var criticalSection bool

func LoadConfiguration(file string) network.Parameters {
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

	demandWait(wait)
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
			Value:      val,
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

func recReceive(processId uint8, req network.RequestCS) {
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

// Main entrypoint for the mutual exclusion program
func main() {
	// Create channels to communicate with the Client Process
	demand := make(chan bool)
	wait := make(chan bool)
	end := make(chan int32) // int32 to pass the new value
	quit := make(chan bool)

	// Create channels to communicate with the Network Process
	message := make(chan []byte)

	// Handle program arguments
	// TODO Possible to move in client??? Or in main.go
	var processId uint8

	if len(os.Args) == 2 {
		tmp, _ := strconv.Atoi(os.Args[1])
		processId = uint8(tmp)
	} else {
		processId = 0
	}

    // Read parameters from json file
    network.Params = LoadConfiguration(parametersFile)

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
			go endDemand(processId, val)

		// Other site releases critical section via Network
		case receivedMsg := <- message:
		    if receivedMsg[0] == byte(network.ReqType) {
		        req := network.DecodeRequest(receivedMsg)
                go recReceive(processId, req)
            } else if receivedMsg[0] == byte(network.OkType) {
                ok := network.DecodeRelease(receivedMsg)
                go okReceive(ok)
            }
            // TODO set value message

        case <-quit:
            return
		}
	}
}
