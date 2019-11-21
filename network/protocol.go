/*
Lab 2 - mutual exclusion
File: network/protocol.go
Authors: Jael Dubey, Luc Wachter
Go version: 1.13.4 (linux/amd64)

Describe networking values, messages structure for the protocol and provide
encoding and decoding functions for messages
*/
package network

import (
    "bytes"
    "encoding/binary"
    "log"
)

// Networking values
// TODO: Maybe parameterize this through json?
const (
    Port = 9706
    Demand = 0
    Wait = 1
    End = 2
)

// Message to request the critical section
type RequestCS struct {
    ProcessNbr uint32
	Timestamp  uint32
}

// Message to release the critical section
type ReleaseCS struct {
    ProcessNbr uint32
	Timestamp  uint32
	Value      int32
}

// Encode given struct as big endian bytes and return bytes buffer
func Encode(message interface{}) []byte {
	buffer := &bytes.Buffer{}
	// Write struct's data as bytes
	err := binary.Write(buffer, binary.BigEndian, message)
	if err != nil {
		log.Fatal(err)
	}

	return buffer.Bytes()
}

// Decode bytes from RequestCS back to struct
func DecodeRequest(buffer []byte) RequestCS {
	message := RequestCS{}
	err := binary.Read(bytes.NewReader(buffer), binary.BigEndian, &message)
	if err != nil {
		log.Fatal(err)
	}

	return message
}

// Decode bytes from ReleaseCS back to struct
func DecodeRelease(buffer []byte) ReleaseCS {
    message := ReleaseCS{}
    err := binary.Read(bytes.NewReader(buffer), binary.BigEndian, &message)
    if err != nil {
        log.Fatal(err)
    }

    return message
}
