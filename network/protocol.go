/*
Lab 2 - mutual exclusion
File: network/protocol.go
Authors: Jael Dubey, Luc Wachter
Go version: 1.13.4 (linux/amd64)

*/
package network

// Networking values
// TODO: Maybe parameterize this through json?
const (
    Port = 9706
    Demand = 0
    Wait = 1
    End = 2
)

type RequestCS struct {
    hi uint32
    i uint32
}

type ReleaseCS struct {
    hi uint32
    i uint32
    value int32
}

