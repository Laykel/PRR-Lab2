/*
Lab 2 - mutual exclusion
File: mutex/mutex.go
Authors: Jael Dubey, Luc Wachter
Go version: 1.13.4 (linux/amd64)

This file contains the implementation of the Carvalho-Roucairol algorithm.
*/
package mutex

import (
	"../client"
	"../network"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

// Necessary algorithm variables
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

func PopulatePWait(processId uint8) {
	// Populate pWait
	for i := processId+1; i < network.Params.NbProcesses; i++ {
			pWait[i] = true
	}
}

// Calculate proper address and forward to network
func send(message []byte, recipientId uint8) {
    recipientPort := strconv.Itoa(int(network.Params.InitialPort + uint16(recipientId)))
    recipientAddress := network.Params.ProcessAddress+":"+recipientPort

    network.Send(message, recipientAddress)
}

func MakeDemand(processId uint8, wait chan bool) {
	timestamp++
	currentDemand = true
	demandTimestamp = timestamp

	// For every process in pWait
	for k := range pWait {
		request := network.MessageCS{
			ReqType:    network.RequestMessageType,
			ProcessNbr: processId,
			Timestamp:  timestamp,
		}

		// Encode message and send to recipient
		send(network.Encode(request), k)
	}

    // Active wait until the pWait list is empty
    for len(pWait) != 0 {
    }

    // Unlock client process
    criticalSection = true
    wait <- true
}

func EndDemand(processId uint8, val int32) {
	timestamp++
	criticalSection = false
	currentDemand = false

	// We send the new value to everybody
	for i := uint8(0); i < network.Params.NbProcesses ; i++ {
		if processId != i {
			val := network.SetVariable{
				ReqType: network.SetValueMessageType,
				Value:   val,
			}
			send(network.Encode(val), i)
		}
	}

	for k := range pDiff {
		ok := network.MessageCS{
			ReqType:    network.ReleaseMessageType,
			ProcessNbr: processId,
			Timestamp:  timestamp,
		}

		// Encode message and send to recipient
		send(network.Encode(ok), k)
	}

	pWait = pDiff
	pDiff = make(map[uint8]bool)
}

func OkReceive(ok network.MessageCS) {
	timestamp = uint32(Max(int64(timestamp), int64(ok.Timestamp)) + 1)

	delete(pWait, ok.ProcessNbr)
}

func ReqReceive(processId uint8, req network.MessageCS) {
	timestamp = uint32(Max(int64(timestamp), int64(req.Timestamp)) + 1)

	if currentDemand == false {
		ok := network.MessageCS{
			ReqType:    network.ReleaseMessageType,
			ProcessNbr: processId,
			Timestamp:  timestamp,
		}

		// Encode message and send to recipient
		send(network.Encode(ok), req.ProcessNbr)
		pWait[req.ProcessNbr] = true
	} else {
		if criticalSection || demandTimestamp < req.Timestamp || ((demandTimestamp == req.Timestamp) && (processId < req.ProcessNbr)) {
			pDiff[req.ProcessNbr] = true
		} else {
			val := network.SetVariable{
				ReqType: network.SetValueMessageType,
				Value:   client.Shared,
			}

			ok := network.MessageCS{
				ReqType:    network.ReleaseMessageType,
				ProcessNbr: processId,
				Timestamp:  timestamp,
			}

			send(network.Encode(val), req.ProcessNbr)
			// Encode message and send to recipient
			send(network.Encode(ok), req.ProcessNbr)
			pWait[req.ProcessNbr] = true

			request := network.MessageCS{
				ReqType:    network.RequestMessageType,
				ProcessNbr: processId,
				Timestamp:  timestamp,
			}

			// Encode message and send to recipient
			send(network.Encode(request), req.ProcessNbr)
		}
	}
}
