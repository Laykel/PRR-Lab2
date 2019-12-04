package network

import (
    "net"
    "testing"
)

var serverAddress = "127.0.0.1:9706"

func init() {
    // Start TCP server
    ch := make(chan []byte)
    go Listen(serverAddress, ch)
}

func TestTCPServerListen(t *testing.T) {
    // Check that server accepts connection
    conn, err := net.Dial("tcp", serverAddress)
    if err != nil {
        t.Error("Error connecting to TCP server: ", err)
    }
    defer conn.Close()
}

func TestTCPServerRequest(t *testing.T) {

}