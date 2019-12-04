package network

import (
    "net"
    "strconv"
    "testing"
)

var serverAddress = "127.0.0.1:"+strconv.Itoa(int(Params.InitialPort))

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
        return
    }
    defer conn.Close()
}

func TestTCPServerRequest(t *testing.T) {

}